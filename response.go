package campaign

import "time"

type Response struct {
	Error      string `json:"error,omitzero"`
	StatusCode int    `json:"status_code,omitzero"`
	Auth       struct {
		Token string `json:"token,omitzero"`
	} `json:"auth,omitzero"`
	Campaign struct {
		ID              uint      `json:"id,omitzero"`
		UserID          string    `json:"user_id,omitzero"`
		Name            string    `json:"name,omitzero"`
		Description     string    `json:"description,omitzero"`
		StartDate       time.Time `json:"start_date,omitzero"`
		EndDate         time.Time `json:"end_date,omitzero"`
		Budget          int64     `json:"budget,omitzero"`
		Status          string    `json:"status,omitzero"`
		ChangeStartDate time.Time `json:"change_start_date,omitzero"`
		ChangeEndDate   time.Time `json:"change_end_date,omitzero"`
		ChangeBudget    int64     `json:"change_budget,omitzero"`
		CreatedAt       time.Time `json:"created_at,omitzero"`
		UpdatedAt       time.Time `json:"updated_at,omitzero"`
	} `json:"campaign,omitzero"`
	Campaigns []struct {
		ID              uint      `json:"id,omitempty"`
		UserID          string    `json:"user_id,omitempty"`
		Name            string    `json:"name,omitempty"`
		Description     string    `json:"description,omitempty"`
		StartDate       time.Time `json:"start_date,omitempty"`
		EndDate         time.Time `json:"end_date,omitempty"`
		Budget          int64     `json:"budget,omitempty"`
		Status          string    `json:"status,omitempty"`
		ChangeStartDate time.Time `json:"change_start_date,omitempty"`
		ChangeEndDate   time.Time `json:"change_end_date,omitempty"`
		ChangeBudget    int64     `json:"change_budget,omitempty"`
		CreatedAt       time.Time `json:"created_at,omitempty"`
		UpdatedAt       time.Time `json:"updated_at,omitempty"`
	} `json:"campaigns,omitempty"`
	Influencer struct {
		InfluencerID string `json:"influencer_id,omitzero"`
		Name         string `json:"name,omitzero"`
	} `json:"influencer,omitzero"`
	Influencers []Influencer `json:"influencers,omitempty"`
	Plan        struct {
		PlanID    string    `json:"plan_id,omitzero"`
		Name      string    `json:"name,omitzero"`
		StartDate time.Time `json:"start_date,omitzero"`
		EndDate   time.Time `json:"end_date,omitzero"`
	} `json:"plan,omitzero"`
}

type Influencer struct {
	InfluencerID string `json:"influencer_id,omitempty"`
	Name         string `json:"name,omitempty"`
}
