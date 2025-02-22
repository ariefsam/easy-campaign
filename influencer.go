package campaign

import "context"

type InfluencerService struct {
	eventStore eventStore
}

func NewInfluencerService() *InfluencerService {
	return &InfluencerService{}
}

func (s *InfluencerService) SetEventStore(eventStore eventStore) {
	s.eventStore = eventStore
}

func (s *InfluencerService) CreateInfluencer(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	return
}
