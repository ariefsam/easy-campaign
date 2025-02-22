package campaign

import (
	"campaign/dto"
	"campaign/eventstore"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
)

type InfluencerService struct {
	eventStore  eventStore
	idGenerator idGenerator
}

func NewInfluencerService() *InfluencerService {
	ctx := context.TODO()
	defaultEventStore, _ := eventstore.New(ctx)
	defaultIDGenerator := idgenerator.New()
	return &InfluencerService{
		eventStore:  defaultEventStore,
		idGenerator: defaultIDGenerator,
	}
}

func (s *InfluencerService) SetEventStore(eventStore eventStore) {
	s.eventStore = eventStore
}

func (s *InfluencerService) SetIDGenerator(idGenerator idGenerator) {
	s.idGenerator = idGenerator
}

func (s *InfluencerService) CreateInfluencer(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	influencerID := s.idGenerator.Generate(ctx)
	event := dto.Event{}
	event.Influencer.InfluencerCreated.InfluencerID = influencerID
	event.Influencer.InfluencerCreated.Name = payload.CreateInfluencerRequest.Name
	event.Influencer.InfluencerCreated.CreatedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Influencer.InfluencerID = influencerID
	resp.Influencer.Name = payload.CreateInfluencerRequest.Name
	return
}
