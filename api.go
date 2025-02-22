package campaign

import (
	"time"
)

type API struct {
	Endpoints []Endpoint `json:"endpoints,omitzero"`
}

type Endpoint struct {
	URL    string   `json:"url,omitzero"`
	Method []string `json:"method,omitzero"`
	Steps  []Step   `json:"handler,omitzero"`
}

type Request struct {
	CreateInfluencerRequest struct {
		Name string `json:"name,omitzero"`
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
}
