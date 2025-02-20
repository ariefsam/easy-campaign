package campaign_test

import (
	"campaign/campaign"
	"campaign/dto"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_campaignService_Create(t *testing.T) {
	ctx := context.TODO()
	campaignService, err := campaign.New()
	require.NoError(t, err)
	require.NotNil(t, campaignService)

	now := time.Date(2025, time.February, 18, 10, 0, 0, 0, time.Local)

	dataEvent := dto.Event{}
	dataEvent.Campaign.CampaignCreated.CampaignID = 1
	dataEvent.Campaign.CampaignCreated.UserID = "userXXX"
	dataEvent.Campaign.CampaignCreated.Name = "campaignXXX"
	dataEvent.Campaign.CampaignCreated.Description = "descriptionXXX"
	dataEvent.Campaign.CampaignCreated.StartDate = now
	dataEvent.Campaign.CampaignCreated.EndDate = now.AddDate(0, 0, 7)
	dataEvent.Campaign.CampaignCreated.Budget = 1000

	err = campaignService.Project(ctx, "e1", dataEvent, now)
	require.NoError(t, err)

	camp, err := campaignService.GetCampaign(ctx, dataEvent.Campaign.CampaignCreated.CampaignID)
	require.NoError(t, err)
	require.NotNil(t, camp)

	require.Equal(t, dataEvent.Campaign.CampaignCreated.CampaignID, camp.ID)
	require.Equal(t, "userXXX", camp.UserID)
	require.Equal(t, "campaignXXX", camp.Name)
	require.Equal(t, "descriptionXXX", camp.Description)
	camp.StartDate = camp.StartDate.In(time.Local)
	require.Equal(t, now, camp.StartDate)
	camp.EndDate = camp.EndDate.In(time.Local)
	require.Equal(t, now.AddDate(0, 0, 7), camp.EndDate)
	require.Equal(t, int64(1000), camp.Budget)

	cursor, err := campaignService.GetCursor()
	require.NoError(t, err)
	require.NotNil(t, cursor)
	require.Equal(t, "e1", cursor.EventID)
}
