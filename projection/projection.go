package projection

import (
	"campaign/dto"
	"campaign/eventstore"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
)

type projection struct {
	projectors  []projector
	redisClient *redis.Client
	eventstore  *eventstore.EventStoreService
}

func New(es *eventstore.EventStoreService) *projection {
	return &projection{
		redisClient: es.GetRedisClient(),
		eventstore:  es,
	}
}

func (p *projection) Run(ctx context.Context, replayFrom string) (err error) {
	if replayFrom != "" {
		logger.Println("replaying from", replayFrom)
		projFuncs := []eventstore.ProjectionFunction{}
		for _, proj := range p.projectors {
			proj.Reset()
			projFuncs = append(projFuncs, proj.Project)
		}
		err = p.eventstore.Replay(ctx, replayFrom, projFuncs)
		if err != nil {
			logger.Error(err)
			return
		}

	}
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
			err = errors.Wrap(resp.Err(), "failed to read group")
			logger.Error(err)
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
	Reset()
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
