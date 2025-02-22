package campaign_test

import (
	"campaign"
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
	resp := &campaign.Response{}
	eventStore := &mockEventStore{}
	influencerService.SetEventStore(eventStore)
	eventStore.On("Save", mock.Anything).Return(nil)

	// payload.
	err := influencerService.CreateInfluencer(ctx, payload, state, resp)
	require.NoError(t, err)
}
