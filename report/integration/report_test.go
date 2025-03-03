package report_test

import (
	"campaign"
	"campaign/eventstore"
	"campaign/logger"
	"campaign/projection"
	"campaign/report"
	"campaign/session"
	"context"
	"log"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {

	log.Default().SetFlags(log.LstdFlags | log.Llongfile)
	ctx, _ := context.WithCancel(context.Background())

	godotenv.Load()
	var err error

	eventstoreService, err := eventstore.New(ctx)
	require.NoError(t, err)

	go func() {
		for {
			eventstoreService.StoreEvent(ctx)
		}
	}()

	rclient := eventstoreService.GetRedisClient()
	require.NotNil(t, rclient)
	rclient.FlushDB(ctx)

	reportService := report.New()
	require.NotNil(t, reportService)

	sessionService, err := session.New()
	require.NoError(t, err)

	projection := projection.New()
	projection.Register(reportService)
	projection.Register(sessionService)
	go projection.Run(ctx)

	influencerService := campaign.NewInfluencerService(eventstoreService)

	payload := &campaign.Request{}
	payload.CreateInfluencerRequest.Name = "influencer1" + time.Now().String()
	influencerName := payload.CreateInfluencerRequest.Name
	state := &campaign.InternalState{}
	resp := &campaign.Response{}
	influencerService.CreateInfluencer(ctx, payload, state, resp)

	authService, err := campaign.NewAuthService()
	require.NoError(t, err)
	payload = &campaign.Request{}
	payload.Login.Email = "admin@gmail.com"
	payload.Login.Password = "makanSaja123!@#"
	state = &campaign.InternalState{}
	resp = &campaign.Response{}
	err = authService.Login(ctx, payload, state, resp)
	require.NoError(t, err)
	token := resp.Auth.Token

	time.Sleep(1 * time.Second)

	sess, err := authService.ParseToken(ctx, token)
	require.NoError(t, err)

	getSession, err := sessionService.GetSession(sess.Id)
	logger.JSON(getSession)
	require.NoError(t, err)
	require.NotNil(t, getSession)
	require.Equal(t, "root", getSession.UserID)
	require.Equal(t, "admin@gmail.com", getSession.Email)

	influencers, err := reportService.FetchInfluencers()
	require.Nil(t, err)
	require.NotZero(t, len(influencers))
	if len(influencers) > 0 {
		require.Equal(t, influencerName, influencers[0].Name)
	}

	influencer, err := reportService.GetInfluencer(influencers[0].InfluencerID)
	require.Nil(t, err)
	require.NotNil(t, influencer)
	require.Equal(t, influencerName, influencer.Name)
	require.NotEmpty(t, influencer.InfluencerID)

	updateRequest := &campaign.Request{}
	updateRequest.UpdateInfluencerRequest.InfluencerID = influencer.InfluencerID
	updateRequest.UpdateInfluencerRequest.Name = "influencer2" + time.Now().String()
	err = influencerService.UpdateInfluencer(ctx, updateRequest, state, resp)
	require.Nil(t, err)
}
