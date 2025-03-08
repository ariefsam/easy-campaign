package campaign

import "time"

type Request struct {
	CreateInfluencerRequest struct {
		Name              string `json:"name,omitzero"`
		TiktokUsername    string `json:"tiktok_username,omitzero"`
		InstagramUsername string `json:"instagram_username,omitzero"`
	} `json:"create_influencer_request,omitzero"`
	Login struct {
		Email       string `json:"email,omitzero"`
		Password    string `json:"password,omitzero"`
		Device      string `json:"device,omitzero"`
		Application string `json:"application,omitzero"` // admin-panel, campaign-dashboard
	} `json:"login,omitzero"`
	CreateCampaignRequest struct {
		UserID      string    `json:"user_id,omitzero"`
		Name        string    `json:"name,omitzero"`
		Description string    `json:"description,omitzero"`
		StartDate   time.Time `json:"start_date,omitzero"`
		EndDate     time.Time `json:"end_date,omitzero"`
		Budget      int64     `json:"budget,omitzero"`
	} `json:"create_campaign_request,omitzero"`
	UpdateCampaignRequest struct {
		ID              uint      `json:"id,omitzero"`
		Name            string    `json:"name,omitzero"`
		Description     string    `json:"description,omitzero"`
		StartDate       time.Time `json:"start_date,omitzero"`
		EndDate         time.Time `json:"end_date,omitzero"`
		Budget          int64     `json:"budget,omitzero"`
		ChangeStartDate time.Time `json:"change_start_date,omitzero"`
		ChangeEndDate   time.Time `json:"change_end_date,omitzero"`
		ChangeBudget    int64     `json:"change_budget,omitzero"`
		Status          string    `json:"status,omitzero"`
	}
	CreatePlanRequest struct {
		Name      string    `json:"name,omitzero"`
		StartDate time.Time `json:"start_date,omitzero"`
		EndDate   time.Time `json:"end_date,omitzero"`
	} `json:"create_plan_request,omitzero"`
	UpdateInfluencerRequest struct {
		InfluencerID string `json:"influencer_id,omitzero"`
		Name         string `json:"name,omitzero"`
	} `json:"update_influencer_request,omitzero"`
	DeleteInfluencerRequest struct {
		InfluencerID string `json:"influencer_id,omitzero"`
	} `json:"delete_influencer_request,omitzero"`
}
