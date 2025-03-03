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

func (m *mockReportService) GetInfluencer(influencerID string) (influencer *report.Influencer, err error) {
	args := m.Called(influencerID)
	influencer, _ = args.Get(0).(*report.Influencer)
	return influencer, args.Error(1)
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

func TestUpdateInfluencerSuccess(t *testing.T) {
	eventStore := &mockEventStore{}
	mockReportService := &mockReportService{}
	influencerService := campaign.NewInfluencerService(eventStore)
	influencerService.SetReportService(mockReportService)

	ctx := context.TODO()
	payload := &campaign.Request{}
	state := &campaign.InternalState{}
	state.Session.UserID = "user123"
	response := &campaign.Response{}

	eventStore.On("Save", mock.Anything).Return(nil)

	payload.UpdateInfluencerRequest.InfluencerID = "ux123"
	payload.UpdateInfluencerRequest.Name = "Ucup Edit"

	expectInfluencer := report.Influencer{
		InfluencerID: "ux123",
		Name:         "Ucup",
	}

	mockReportService.On("GetInfluencer", "ux123").Return(&expectInfluencer, nil)

	expectEvent := dto.Event{}
	expectEvent.Influencer.InfluencerUpdated.Name = payload.UpdateInfluencerRequest.Name
	expectEvent.Influencer.InfluencerUpdated.InfluencerID = payload.UpdateInfluencerRequest.InfluencerID
	expectEvent.Influencer.InfluencerUpdated.UpdatedBy = state.Session.UserID

	err := influencerService.UpdateInfluencer(ctx, payload, state, response)
	require.NoError(t, err)

	eventStore.AssertCalled(t, "Save", expectEvent)
	require.Equal(t, expectInfluencer.InfluencerID, response.Influencer.InfluencerID)
	require.Equal(t, "Ucup Edit", response.Influencer.Name)
}

func TestUpdateInfluencerFailedNotExist(t *testing.T) {
	eventStore := &mockEventStore{}
	mockReport := &mockReportService{}
	influencerService := campaign.NewInfluencerService(eventStore)
	influencerService.SetReportService(mockReport)

	ctx := context.TODO()
	payload := &campaign.Request{}
	payload.UpdateInfluencerRequest.InfluencerID = "ux123"
	payload.UpdateInfluencerRequest.Name = "Ucup Edit"
	state := &campaign.InternalState{}
	response := &campaign.Response{}

	mockReport.On("GetInfluencer", payload.UpdateInfluencerRequest.InfluencerID).Return(nil, nil)

	expectEvent := dto.Event{}
	expectEvent.Influencer.InfluencerUpdated.Name = payload.UpdateInfluencerRequest.Name
	expectEvent.Influencer.InfluencerUpdated.InfluencerID = payload.UpdateInfluencerRequest.InfluencerID
	expectEvent.Influencer.InfluencerUpdated.UpdatedBy = state.Session.UserID

	err := influencerService.UpdateInfluencer(ctx, payload, state, response)
	require.Error(t, err)

	mockReport.AssertCalled(t, "GetInfluencer", "ux123")
	eventStore.AssertNotCalled(t, "Save", mock.Anything)
}
