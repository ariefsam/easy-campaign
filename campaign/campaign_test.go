package campaign_test

import (
	"campaign/campaign"
	"campaign/dto"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// go test -coverprofile=coverage.out ./campaign/...
// go tool cover -func=coverage.out
// go tool cover -html=coverage.out

func Test_campaignService_Create(t *testing.T) {
	ctx := context.TODO()
	campaignService, err := campaign.New()
	require.NoError(t, err)
	require.NotNil(t, campaignService)

	now := time.Date(2025, time.February, 18, 10, 0, 0, 0, time.Local)

	dataEvent := dto.Event{}
	dataEvent.Campaign.CampaignCreated.ID = 1
	dataEvent.Campaign.CampaignCreated.UserID = "userXXX"
	dataEvent.Campaign.CampaignCreated.Name = "campaignXXX"
	dataEvent.Campaign.CampaignCreated.Description = "descriptionXXX"
	dataEvent.Campaign.CampaignCreated.StartDate = now
	dataEvent.Campaign.CampaignCreated.EndDate = now.AddDate(0, 0, 7)
	dataEvent.Campaign.CampaignCreated.Budget = 1000

	err = campaignService.Project(ctx, "e1", dataEvent, now)
	require.NoError(t, err)

	camp, err := campaignService.GetCampaign(ctx, dataEvent.Campaign.CampaignCreated.ID)
	require.NoError(t, err)
	require.NotNil(t, camp)

	require.Equal(t, dataEvent.Campaign.CampaignCreated.ID, camp.ID)
	require.Equal(t, "userXXX", camp.UserID)
	require.Equal(t, "campaignXXX", camp.Name)
	require.Equal(t, "descriptionXXX", camp.Description)
	camp.StartDate = camp.StartDate.In(time.Local)
	require.Equal(t, now, camp.StartDate)
	camp.EndDate = camp.EndDate.In(time.Local)
	require.Equal(t, now.AddDate(0, 0, 7), camp.EndDate)
	require.Equal(t, int64(1000), camp.Budget)
	require.Equal(t, "active", camp.Status)

	dataEventUpdate := dto.Event{}
	dataEventUpdate.Campaign.CampaignUpdated.ID = 1
	dataEventUpdate.Campaign.CampaignUpdated.UserID = "userYYY"
	dataEventUpdate.Campaign.CampaignUpdated.Name = "campaignYYY"
	dataEventUpdate.Campaign.CampaignUpdated.Description = "descriptionYYY"
	dataEventUpdate.Campaign.CampaignUpdated.StartDate = now
	dataEventUpdate.Campaign.CampaignUpdated.EndDate = now.AddDate(0, 0, 5)
	dataEventUpdate.Campaign.CampaignUpdated.Budget = 2000
	dataEventUpdate.Campaign.CampaignUpdated.ChangeStartDate = now
	dataEventUpdate.Campaign.CampaignUpdated.ChangeEndDate = now.AddDate(0, 0, 6)
	dataEventUpdate.Campaign.CampaignUpdated.ChangeBudget = 3000
	dataEventUpdate.Campaign.CampaignUpdated.Status = "inactive"

	err = campaignService.Project(ctx, "e2", dataEventUpdate, now)
	require.NoError(t, err)

	campUpdate, err := campaignService.GetCampaign(ctx, dataEventUpdate.Campaign.CampaignUpdated.ID)
	require.NoError(t, err)
	require.NotNil(t, camp)

	require.Equal(t, dataEventUpdate.Campaign.CampaignUpdated.ID, campUpdate.ID)
	require.Equal(t, "userYYY", campUpdate.UserID)
	require.Equal(t, "campaignYYY", campUpdate.Name)
	require.Equal(t, "descriptionYYY", campUpdate.Description)
	campUpdate.StartDate = campUpdate.StartDate.In(time.Local)
	require.Equal(t, now, campUpdate.StartDate)
	campUpdate.EndDate = campUpdate.EndDate.In(time.Local)
	require.Equal(t, now.AddDate(0, 0, 5), campUpdate.EndDate)
	require.Equal(t, int64(2000), campUpdate.Budget)
	require.Equal(t, "inactive", campUpdate.Status)
	campUpdate.ChangeStartDate = campUpdate.ChangeStartDate.In(time.Local)
	require.Equal(t, now, campUpdate.ChangeStartDate)
	campUpdate.ChangeEndDate = campUpdate.ChangeEndDate.In(time.Local)
	require.Equal(t, now.AddDate(0, 0, 6), campUpdate.ChangeEndDate)
	require.Equal(t, int64(3000), campUpdate.ChangeBudget)

	dataEventDelete := dto.Event{}
	dataEventDelete.Campaign.CampaignDeleted.ID = 1

	err = campaignService.Project(ctx, "e3", dataEventDelete, now)
	require.NoError(t, err)

	cursor, err := campaignService.GetCursor()
	require.NoError(t, err)
	require.NotNil(t, cursor)
	require.Equal(t, "e3", cursor.EventID)

	os.Remove("campaign.db")
}
