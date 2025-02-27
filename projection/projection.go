package projection

import (
	"campaign/dto"
	"campaign/eventstore"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
	"encoding/json"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type projection struct {
	projectors  []projector
	redisClient *redis.Client
}

type eventstoreService interface {
	GetRedisClient() *redis.Client
}

func New() *projection {
	es, err := eventstore.New(context.TODO())
	if err != nil {
		panic(err)
	}
	return &projection{
		redisClient: es.GetRedisClient(),
	}
}

func (p *projection) Run(ctx context.Context) (err error) {
	for _, proj := range p.projectors {
		groupName := proj.GetGroupName()
		keys := proj.SubscribedTo()
		cursor, _ := proj.GetCursor()
		if cursor == "" {
			cursor = "0"
		}
		for _, key := range keys {
			p.RegisterGroup(ctx, key, groupName, cursor)
		}

		go p.project(ctx, proj)

	}
	return
}

func (es *projection) project(ctx context.Context, proj projector) (err error) {
	entities := proj.SubscribedTo()
	streams := []string{}
	groupName := proj.GetGroupName()
	streams = append(streams, entities...)
	consumer := idgenerator.New().Generate(ctx)

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
			if resp.Err() == redis.Nil {
				return
			}
			logger.Println("failed to read group: ", resp.Err())
			time.Sleep(1 * time.Second)
			continue
		}

		streams := resp.Val()

		for _, stream := range streams {
			streamName := stream.Stream
			messages := stream.Messages
			for _, message := range messages {
				dateTime, _ := time.Parse(time.RFC3339, message.Values["time"].(string))
				data := dto.Event{}
				dataRaw, _ := message.Values["data"].(string)
				json.Unmarshal([]byte(dataRaw), &data)

				proj.Project(ctx, message.ID, data, dateTime)
				es.redisClient.XAck(ctx2, streamName, groupName, message.ID)
			}
		}

	}
}

type projector interface {
	Project(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error)
	SubscribedTo() []string
	GetCursor() (eventID string, err error)
	GetGroupName() string
}

func (p *projection) Register(proj projector) {
	p.projectors = append(p.projectors, proj)
}

func (proj *projection) RegisterGroup(ctx context.Context, key, groupName, fromID string) (err error) {

	proj.redisClient.XGroupCreateMkStream(ctx, key, groupName, "0")

	resp := proj.redisClient.XInfoConsumers(ctx, key, groupName)
	consumers := resp.Val()
	if consumers != nil {
		return
	}
	logger.Println("creating group", key, groupName, fromID)
	res := proj.redisClient.XGroupCreate(ctx, key, groupName, fromID)
	err = res.Err()
	if err != nil {
		logger.Println("failed to create group: " + err.Error())
		return
	}

	return
}
