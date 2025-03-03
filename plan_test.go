package campaign_test

import (
	"campaign"
	"campaign/dto"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewPlanService(t *testing.T) {
	ctx := context.TODO()
	es := &mockEventStore{}
	idGenerator := &mockIDGenerator{}
	planService := campaign.NewPlanService(es)
	require.NotNil(t, planService)

	planService.SetIDGenerator(idGenerator)

	payload := &campaign.Request{}
	payload.CreatePlanRequest.Name = "Plan 1"

	state := &campaign.InternalState{}
	resp := &campaign.Response{}

	start, _ := time.Parse(time.RFC3339, "2020-06-18T17:24:53Z")
	end, _ := time.Parse(time.RFC3339, "2020-07-18T17:24:53Z")

	payload.CreatePlanRequest.StartDate = start
	payload.CreatePlanRequest.EndDate = end
	expectedEvent := dto.Event{
		Plan: dto.Plan{
			PlanCreated: dto.PlanCreated{
				PlanID:    "plan123",
				Name:      "Plan 1",
				StartDate: start,
				EndDate:   end,
				CreatedBy: state.Session.UserID,
			},
		},
	}

	idGenerator.On("Generate").Return("plan123").Times(1)
	es.On("Save", mock.Anything).Return(nil)

	err := planService.Create(ctx, payload, state, resp)
	require.NoError(t, err)

	es.AssertCalled(t, "Save", expectedEvent)

}
