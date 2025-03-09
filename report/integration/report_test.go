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
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestReport(t *testing.T) {

	waitDuration := 300 * time.Millisecond

	log.Default().SetFlags(log.LstdFlags | log.Llongfile)
	ctx, _ := context.WithCancel(context.Background())

	func() {
		os.Remove("report.db")
		os.Remove("event.db")
		os.Remove("session.db")
	}()

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

	projection := projection.New(eventstoreService)
	projection.Register(reportService)
	projection.Register(sessionService)
	go projection.Run(ctx, "")

	influencerService := campaign.NewInfluencerService(eventstoreService)

	payload := &campaign.Request{}
	payload.CreateInfluencerRequest.Name = "influencer1" + time.Now().String()
	influencerName := payload.CreateInfluencerRequest.Name
	state := &campaign.InternalState{}
	resp := &campaign.Response{}
	influencerService.CreateInfluencer(ctx, payload, state, resp)
	influencerService.SetReportService(reportService)

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

	time.Sleep(waitDuration)

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

	time.Sleep(waitDuration)
	influencerCheck, err := reportService.GetInfluencer(influencers[0].InfluencerID)
	require.Nil(t, err)
	require.NotNil(t, influencerCheck)
	require.Equal(t, updateRequest.UpdateInfluencerRequest.Name, influencerCheck.Name)

	t.Run("delete influencer", func(t *testing.T) {
		deleteRequest := &campaign.Request{}
		deleteRequest.DeleteInfluencerRequest.InfluencerID = influencer.InfluencerID
		err = influencerService.DeleteInfluencer(ctx, deleteRequest, state, resp)
		require.Nil(t, err)

		time.Sleep(waitDuration)

		influencer, err := reportService.GetInfluencer(influencers[0].InfluencerID)
		require.Nil(t, err)
		require.Nil(t, influencer)

	})

	t.Run("create plan", func(t *testing.T) {
		planService := campaign.NewPlanService(eventstoreService)
		planService.SetPlanProjection(reportService)

		payload := &campaign.Request{}
		payload.CreatePlanRequest.Name = "plan1" + time.Now().String()
		payload.CreatePlanRequest.StartDate = time.Now()
		payload.CreatePlanRequest.EndDate = time.Now().Add(100 * time.Hour)
		state := &campaign.InternalState{}
		resp := &campaign.Response{}

		err := planService.Create(ctx, payload, state, resp)
		require.Nil(t, err)

		time.Sleep(waitDuration)

		plans, err := reportService.FetchPlans()
		require.Nil(t, err)
		require.NotZero(t, len(plans))
		if len(plans) > 0 {
			require.Equal(t, payload.CreatePlanRequest.Name, plans[0].Name)
		}

		plan, err := reportService.GetPlan(plans[0].PlanID)
		require.Nil(t, err)
		require.NotNil(t, plan)
		require.Equal(t, payload.CreatePlanRequest.Name, plan.Name)
		require.NotEmpty(t, plan.PlanID)

		t.Run("update plan", func(t *testing.T) {
			updateRequest := &campaign.Request{}
			updateRequest.UpdatePlanRequest.PlanID = plan.PlanID
			updateRequest.UpdatePlanRequest.Name = "plan2" + time.Now().String()
			updateRequest.UpdatePlanRequest.StartDate = time.Now()
			updateRequest.UpdatePlanRequest.EndDate = time.Now().Add(200 * time.Hour)
			err = planService.Update(ctx, updateRequest, state, resp)
			require.Nil(t, err)

			time.Sleep(waitDuration)

			planCheck, err := reportService.GetPlan(plans[0].PlanID)
			require.Nil(t, err)
			require.NotNil(t, planCheck)
			require.Equal(t, updateRequest.UpdatePlanRequest.Name, planCheck.Name)

		})

		t.Run("delete plan", func(t *testing.T) {
			deleteRequest := &campaign.Request{}
			deleteRequest.DeletePlanRequest.PlanID = plan.PlanID
			err = planService.Delete(ctx, deleteRequest, state, resp)
			require.Nil(t, err)

			time.Sleep(waitDuration)

			plan, err := reportService.GetPlan(plans[0].PlanID)
			require.Nil(t, err)
			require.Nil(t, plan)
		})
	})
}
