package campaign

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/logger"
	"context"
)

type planService struct {
	eventStore  eventStore
	idGenerator idGenerator
}

func NewPlanService(es eventStore) *planService {
	defaultIDGenerator := idgenerator.New()
	return &planService{
		eventStore:  es,
		idGenerator: defaultIDGenerator,
	}
}

func (s *planService) SetIDGenerator(idGenerator idGenerator) {
	s.idGenerator = idGenerator
}

func (s *planService) Create(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	planID := s.idGenerator.Generate(ctx)
	event := dto.Event{}
	event.Plan.PlanCreated.PlanID = planID
	event.Plan.PlanCreated.Name = payload.CreatePlanRequest.Name
	event.Plan.PlanCreated.StartDate = payload.CreatePlanRequest.StartDate
	event.Plan.PlanCreated.EndDate = payload.CreatePlanRequest.EndDate
	event.Plan.PlanCreated.CreatedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Plan.PlanID = planID
	resp.Plan.Name = payload.CreatePlanRequest.Name
	resp.Plan.StartDate = payload.CreatePlanRequest.StartDate
	resp.Plan.EndDate = payload.CreatePlanRequest.EndDate

	return
}
