package campaign_test

import (
	"campaign"
	"campaign/dto"
	"campaign/report"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockPlanProjection struct {
	mock.Mock
}

func (m *mockPlanProjection) GetPlan(planID string) (plan *report.Plan, err error) {
	args := m.Called(planID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*report.Plan), args.Error(1)
}

func TestNewPlanService(t *testing.T) {
	ctx := context.TODO()
	es := &mockEventStore{}
	idGenerator := &mockIDGenerator{}
	planService := campaign.NewPlanService(es)
	require.NotNil(t, planService)

	planService.SetIDGenerator(idGenerator)
	planProjection := &mockPlanProjection{}
	planService.SetPlanProjection(planProjection)

	t.Run("Create Plan", func(t *testing.T) {
		payload := &campaign.Request{}
		payload.CreatePlanRequest.Name = "Plan 1"

		state := &campaign.InternalState{}
		resp := &campaign.Response{}

		start, _ := time.Parse(time.RFC3339, "2020-06-18T17:24:53Z")
		end, _ := time.Parse(time.RFC3339, "2020-07-18T17:24:53Z")

		payload.CreatePlanRequest.StartDate = start
		payload.CreatePlanRequest.EndDate = end
		expectedEvent := dto.Event{}

		expectedEvent.Plan.PlanCreated.PlanID = "plan123"
		expectedEvent.Plan.PlanCreated.Name = "Plan 1"
		expectedEvent.Plan.PlanCreated.StartDate = start
		expectedEvent.Plan.PlanCreated.EndDate = end
		expectedEvent.Plan.PlanCreated.CreatedBy = state.Session.UserID

		idGenerator.On("Generate").Return("plan123").Times(1)
		es.On("Save", mock.Anything).Return(nil)

		err := planService.Create(ctx, payload, state, resp)
		require.NoError(t, err)

		es.AssertCalled(t, "Save", expectedEvent)
	})

	t.Run("Update plan not found", func(t *testing.T) {
		payload := &campaign.Request{}
		payload.UpdatePlanRequest.PlanID = "plan123"
		payload.UpdatePlanRequest.Name = "Plan 1"

		state := &campaign.InternalState{}
		state.Session.UserID = "user123"
		resp := &campaign.Response{}

		planProjection.On("GetPlan", "plan123").Return(nil, nil).Once()

		err := planService.Update(ctx, payload, state, resp)
		require.Error(t, err)

	})

	t.Run("Update plan success", func(t *testing.T) {
		payload := &campaign.Request{}
		payload.UpdatePlanRequest.PlanID = "plan123"
		payload.UpdatePlanRequest.Name = "Plan 1"
		payload.UpdatePlanRequest.StartDate = time.Now()
		payload.UpdatePlanRequest.EndDate = time.Now().Add(300 * time.Hour)

		state := &campaign.InternalState{}
		state.Session.UserID = "user123"
		resp := &campaign.Response{}

		plan := &report.Plan{
			PlanID:    "plan123",
			Name:      "Plan 1",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(100 * time.Hour),
			CreatedBy: "user123",
		}

		planProjection.On("GetPlan", "plan123").Return(plan, nil).Once()

		expectedEvent := dto.Event{}
		expectedEvent.Plan.PlanUpdated.PlanID = "plan123"
		expectedEvent.Plan.PlanUpdated.Name = "Plan 1"
		expectedEvent.Plan.PlanUpdated.StartDate = payload.UpdatePlanRequest.StartDate
		expectedEvent.Plan.PlanUpdated.EndDate = payload.UpdatePlanRequest.EndDate
		expectedEvent.Plan.PlanUpdated.UpdatedBy = state.Session.UserID

		es.On("Save", mock.Anything).Return(nil)

		err := planService.Update(ctx, payload, state, resp)
		require.NoError(t, err)

		es.AssertCalled(t, "Save", expectedEvent)
	})

	t.Run("Delete plan not found", func(t *testing.T) {
		payload := &campaign.Request{}
		payload.DeletePlanRequest.PlanID = "plan123"

		state := &campaign.InternalState{}
		state.Session.UserID = "user123"
		resp := &campaign.Response{}

		planProjection.On("GetPlan", "plan123").Return(nil, nil).Once()

		err := planService.Delete(ctx, payload, state, resp)
		require.Error(t, err)

	})

	t.Run("Delete plan success", func(t *testing.T) {
		payload := &campaign.Request{}
		payload.DeletePlanRequest.PlanID = "plan123"

		state := &campaign.InternalState{}
		state.Session.UserID = "user123"
		resp := &campaign.Response{}

		plan := &report.Plan{
			PlanID:    "plan123",
			Name:      "Plan 1",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(100 * time.Hour),
			CreatedBy: "user123",
		}

		planProjection.On("GetPlan", "plan123").Return(plan, nil).Once()

		expectedEvent := dto.Event{}
		expectedEvent.Plan.PlanDeleted.PlanID = "plan123"
		expectedEvent.Plan.PlanDeleted.DeletedBy = state.Session.UserID

		es.On("Save", mock.Anything).Return(nil)

		err := planService.Delete(ctx, payload, state, resp)
		require.NoError(t, err)

		es.AssertCalled(t, "Save", expectedEvent)
	})

}
