package campaign

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/logger"
	"campaign/report"
	"context"
	"errors"
)

type planProjection interface {
	GetPlan(planID string) (plan *report.Plan, err error)
}
type planService struct {
	eventStore     eventStore
	idGenerator    idGenerator
	planProjection planProjection
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

func (s *planService) SetPlanProjection(planProjection planProjection) {
	s.planProjection = planProjection
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

func (s *planService) Update(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	plan, err := s.planProjection.GetPlan(payload.UpdatePlanRequest.PlanID)
	if err != nil {
		logger.Error(err)
		return
	}

	if plan == nil {
		err = errors.New("plan not found: " + payload.UpdatePlanRequest.PlanID)
		logger.Error(err)
		return
	}

	event := dto.Event{}
	event.Plan.PlanUpdated.PlanID = plan.PlanID
	event.Plan.PlanUpdated.Name = payload.UpdatePlanRequest.Name
	event.Plan.PlanUpdated.StartDate = payload.UpdatePlanRequest.StartDate
	event.Plan.PlanUpdated.EndDate = payload.UpdatePlanRequest.EndDate
	event.Plan.PlanUpdated.UpdatedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Plan.PlanID = plan.PlanID
	resp.Plan.Name = payload.UpdatePlanRequest.Name
	resp.Plan.StartDate = payload.UpdatePlanRequest.StartDate
	resp.Plan.EndDate = payload.UpdatePlanRequest.EndDate

	return
}

func (s *planService) Delete(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	plan, err := s.planProjection.GetPlan(payload.DeletePlanRequest.PlanID)
	if err != nil {
		logger.Error(err)
		return
	}

	if plan == nil {
		err = errors.New("plan not found: " + payload.DeletePlanRequest.PlanID)
		logger.Error(err)
		return
	}

	event := dto.Event{}
	event.Plan.PlanDeleted.PlanID = plan.PlanID
	event.Plan.PlanDeleted.DeletedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Plan.PlanID = plan.PlanID
	resp.Plan.Name = plan.Name
	resp.Plan.StartDate = plan.StartDate
	resp.Plan.EndDate = plan.EndDate

	return
}
