package eventstore

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type eventStoreService struct {
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

func New(ctx context.Context) (es *eventStoreService, err error) {
	client, err := gorm.Open(sqlite.Open("event.db"), &gorm.Config{})
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

	es = &eventStoreService{
		db:          client,
		redisClient: rdb,
	}

	err = es.db.AutoMigrate(&Event{})
	if err != nil {
		err = errors.Wrap(err, "failed to migrate event")
		logger.Println(err)
	}

	go es.StoreEvent(ctx)

	return
}

func (es *eventStoreService) Save(ctx context.Context, eventData dto.Event) (err error) {
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

	logger.PrintJSON(string(jsonData))
	return
}

func (es *eventStoreService) StoreEvent(ctx context.Context) {
	idService := idgenerator.New()
	consumer := idService.Generate(ctx)
	groupName := "eventstore"
	entities := dto.ExtractListEntities(dto.Event{})

	lastEvent := Event{}

	es.db.Debug().Last(&lastEvent)
	if lastEvent.EventID == "" {
		lastEvent.EventID = "0"
	}

	for _, entity := range entities {
		err := es.RegisterGroup(ctx, entity, groupName, lastEvent.EventID)
		if err != nil {
			err = errors.Wrap(err, "failed to register group")
			logger.Println(err)
			panic(err)
		}
	}

	streams := []string{}

	streams = append(streams, entities...)

	for _ = range entities {
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
			logger.Println("failed to read group: ", resp.Err())
			time.Sleep(1 * time.Second)
			continue
		}

		streams := resp.Val()

		for _, stream := range streams {
			streamName := stream.Stream
			messages := stream.Messages

			for _, message := range messages {
				logger.PrintJSON(message.Values)
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
					logger.Println(err)
					continue
				}

				logger.PrintJSON(event)
				es.redisClient.XAck(ctx2, streamName, groupName, message.ID)

			}
		}

	}

}

func (es *eventStoreService) RegisterGroup(ctx context.Context, key, groupName, fromID string) (err error) {

	es.redisClient.XGroupCreateMkStream(ctx, key, groupName, "0")

	resp := es.redisClient.XInfoConsumers(ctx, key, groupName)
	consumers := resp.Val()
	if consumers != nil {
		return
	}
	logger.Println("creating group", key, groupName, fromID)
	res := es.redisClient.XGroupCreate(ctx, key, groupName, fromID)
	err = res.Err()
	if err != nil {
		logger.Println("failed to create group: " + err.Error())
		return
	}

	return
}
