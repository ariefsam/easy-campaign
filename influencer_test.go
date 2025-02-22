package campaign_test

import (
	"campaign"
	"campaign/dto"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInfluencerService_CreateInfluencer(t *testing.T) {
	influencerService := campaign.NewInfluencerService()

	ctx := context.TODO()
	payload := &campaign.Request{}
	state := &campaign.InternalState{}
	state.Session.UserID = "user123"
	resp := &campaign.Response{}
	eventStore := &mockEventStore{}
	idGenerator := &mockIDGenerator{}
	influencerService.SetEventStore(eventStore)
	influencerService.SetIDGenerator(idGenerator)
	expectID := "inf123"
	idGenerator.On("Generate").Return(expectID)

	eventStore.On("Save", mock.Anything).Return(nil)

	payload.CreateInfluencerRequest.Name = "Ucup"

	expectEvent := dto.Event{}
	expectEvent.Influencer.InfluencerCreated.Name = "Ucup"
	expectEvent.Influencer.InfluencerCreated.InfluencerID = expectID
	expectEvent.Influencer.InfluencerCreated.CreatedBy = state.Session.UserID

	err := influencerService.CreateInfluencer(ctx, payload, state, resp)
	require.NoError(t, err)

	eventStore.AssertCalled(t, "Save", expectEvent)
	require.Equal(t, expectID, resp.Influencer.InfluencerID)
	require.Equal(t, "Ucup", resp.Influencer.Name)

}
