package eventstore

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type EventStoreService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

type Event struct {
	gorm.Model
	EventID  string    `json:"event_id" gorm:"unique"`
	Entity   string    `json:"entity" gorm:"index"`
	Event    string    `json:"event" gorm:"index"`
	Data     string    `json:"data"`
	DateTime time.Time `json:"datetime"`
}

func New(ctx context.Context) (es *EventStoreService, err error) {
	client, err := gorm.Open(sqlite.Open("event.db"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		err = errors.Wrap(err, "failed to connect to database")
		logger.Println(err)
		return
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		err = errors.New("redis host is empty")
		logger.Println(err)
		return
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	es = &EventStoreService{
		db:          client,
		redisClient: rdb,
	}

	err = es.db.AutoMigrate(&Event{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate event")
		logger.Println(err)
	}

	return
}

func (es *EventStoreService) Save(ctx context.Context, eventData dto.Event) (err error) {
	jsonData, err := json.Marshal(eventData)
	if err != nil {
		err = errors.Wrap(err, "failed to marshal event")
		logger.Println(err)
		return
	}

	entityName, event := dto.ExtractEvent(eventData)

	now := time.Now()

	xaddCmd := es.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: entityName,
		Values: map[string]interface{}{
			"data":  string(jsonData),
			"event": event,
			"time":  now,
		},
	})

	if xaddCmd.Err() != nil {
		err = errors.Wrap(xaddCmd.Err(), "failed to add event")
		logger.Println(err)
		return
	}

	return
}

func (es *EventStoreService) GetRedisClient() *redis.Client {
	return es.redisClient
}

func (es *EventStoreService) StoreEvent(ctx context.Context) {
	idService := idgenerator.New()
	consumer := idService.Generate(ctx)
	groupName := "eventstore"
	ev := dto.Event{}
	entities := ev.GetEntityList()

	lastEvent := Event{}

	es.db.Last(&lastEvent)
	if lastEvent.EventID == "" {
		lastEvent.EventID = "0"
	}

	for _, entity := range entities {
		err := es.RegisterGroup(ctx, entity, groupName, lastEvent.EventID)
		if err != nil {
			err = errors.Wrap(err, "failed to register group")
			logger.Println(err)
		}
	}

	streams := []string{}

	streams = append(streams, entities...)

	for range entities {
		streams = append(streams, ">")
	}

	for {
		if ctx.Err() != nil {
			logger.PrintJSON(ctx.Err())
			return
		}

		ctx2 := context.Background()
		resp := es.redisClient.XReadGroup(ctx2, &redis.XReadGroupArgs{
			Group:    groupName,
			Streams:  streams,
			Consumer: consumer,
			Count:    1,
			Block:    10000 * time.Millisecond,
		})

		if resp.Err() != nil {
			if resp.Err() == redis.Nil {
				return
			}
			logger.Println("failed to read group: ", resp.Err())
			time.Sleep(1 * time.Second)
			sentry.CaptureException(resp.Err())
			continue
		}

		streams := resp.Val()

		for _, stream := range streams {
			streamName := stream.Stream
			messages := stream.Messages

			for _, message := range messages {
				dateTime, _ := time.Parse(time.RFC3339, message.Values["time"].(string))
				event := Event{
					EventID:  message.ID,
					Entity:   streamName,
					Event:    message.Values["event"].(string),
					Data:     message.Values["data"].(string),
					DateTime: dateTime,
				}

				err := es.db.Create(&event).Error
				if err != nil {
					err = errors.Wrap(err, "failed to create event")
					logger.Error(err)
					continue
				}

				es.redisClient.XAck(ctx2, streamName, groupName, message.ID)

			}
		}

	}

}

func (es *EventStoreService) RegisterGroup(ctx context.Context, key, groupName, fromID string) (err error) {

	crtStream := es.redisClient.XGroupCreateMkStream(ctx, key, groupName, "0")
	if crtStream.Err() != nil {
		if crtStream.Err().Error() == "BUSYGROUP Consumer Group name already exists" {
			return
		}
		logger.Println("key:", key, "groupName:", groupName, "fromID:", fromID, ".failed to create stream:"+crtStream.Err().Error())
		return
	}

	resp := es.redisClient.XInfoConsumers(ctx, key, groupName)
	consumers := resp.Val()
	if consumers != nil {
		return
	}
	logger.Println("creating group", key, groupName, fromID)
	res := es.redisClient.XGroupCreate(ctx, key, groupName, fromID)
	err = res.Err()
	if err != nil {
		logger.Println("key:", key, "groupName:", groupName, "fromID:", fromID, ". failed to create group: "+err.Error())
		return
	}

	return
}

type ProjectionFunction func(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) error

func (es *EventStoreService) Replay(ctx context.Context, fromEvent string, projectionFunc []ProjectionFunction) (err error) {
	event := Event{}
	if fromEvent == "" || fromEvent == "0" {
		err = es.db.First(&event).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(err, "failed to get first event")
			logger.Error(err)
			return
		}
	} else {
		err = es.db.First(&event, "event_id=?", fromEvent).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(err, "failed to get event")
			logger.Error(err)
			return
		}
	}

	logger.Println("replaying from event:", event.EventID)

	rows, err := es.db.Model(&Event{}).Where("id > ?", event.ID).Rows()
	if err != nil {
		err = errors.Wrap(err, "failed to get rows")
		logger.Error(err)
		return
	}

	for rows.Next() {
		event := Event{}
		err = es.db.ScanRows(rows, &event)
		if err != nil {
			err = errors.Wrap(err, "failed to scan rows")
			logger.Error(err)
			continue
		}

		logger.PrintJSON("replaying event:")
		logger.PrintJSON(event)

		for _, proj := range projectionFunc {
			dtoEvent := dto.Event{}
			err = json.Unmarshal([]byte(event.Data), &dtoEvent)
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal event")
				logger.Error(err)
				continue
			}

			err = proj(ctx, event.EventID, dtoEvent, event.DateTime)
			if err != nil {
				err = errors.Wrap(err, "failed to project")
				logger.Error(err)
				continue
			}
		}

	}

	return
}
