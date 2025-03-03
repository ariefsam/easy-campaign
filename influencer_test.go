package campaign_test

import (
	"campaign"
	"campaign/dto"
	"campaign/report"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInfluencerService_CreateInfluencer(t *testing.T) {
	eventStore := &mockEventStore{}
	influencerService := campaign.NewInfluencerService(eventStore)

	ctx := context.TODO()
	payload := &campaign.Request{}
	state := &campaign.InternalState{}
	state.Session.UserID = "user123"
	resp := &campaign.Response{}

	idGenerator := &mockIDGenerator{}
	influencerService.SetIDGenerator(idGenerator)
	expectID := "inf123"
	idGenerator.On("Generate").Return(expectID)

	eventStore.On("Save", mock.Anything).Return(nil)

	payload.CreateInfluencerRequest.Name = "Ucup"
	payload.CreateInfluencerRequest.InstagramUsername = "ucupgram"
	payload.CreateInfluencerRequest.TiktokUsername = "ucuptok"

	expectEvent := dto.Event{}
	expectEvent.Influencer.InfluencerCreated.Name = "Ucup"
	expectEvent.Influencer.InfluencerCreated.InstagramUsername = "ucupgram"
	expectEvent.Influencer.InfluencerCreated.TiktokUsername = "ucuptok"
	expectEvent.Influencer.InfluencerCreated.InfluencerID = expectID
	expectEvent.Influencer.InfluencerCreated.CreatedBy = state.Session.UserID

	err := influencerService.CreateInfluencer(ctx, payload, state, resp)
	require.NoError(t, err)

	eventStore.AssertCalled(t, "Save", expectEvent)
	require.Equal(t, expectID, resp.Influencer.InfluencerID)
	require.Equal(t, "Ucup", resp.Influencer.Name)

}

type mockReportService struct {
	mock.Mock
}

func (m *mockReportService) FetchInfluencers() (influencers []report.Influencer, err error) {
	args := m.Called()
	return args.Get(0).([]report.Influencer), args.Error(1)
}

func TestFetchInfluencer(t *testing.T) {
	eventStore := &mockEventStore{}
	influencerService := campaign.NewInfluencerService(eventStore)

	ctx := context.TODO()
	payload := &campaign.Request{}
	state := &campaign.InternalState{}
	response := &campaign.Response{}

	reportService := &mockReportService{}

	t.Run("report service is not set", func(t *testing.T) {
		err := influencerService.FetchInfluencers(ctx, payload, state, response)
		require.Error(t, err)
	})

	influencerService.SetReportService(reportService)

	t.Run("fetch influencers error", func(t *testing.T) {

		reportService.On("FetchInfluencers").Return([]report.Influencer{
			{
				InfluencerID: "inf123",
				Name:         "Ucup",
			},
		}, nil)

		err := influencerService.FetchInfluencers(ctx, payload, state, response)
		require.NoError(t, err)

		reportService.AssertCalled(t, "FetchInfluencers")
	})

}
