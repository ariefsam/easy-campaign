package campaign

import (
	"context"
	"time"
)

type Application struct {
	Event Event `json:"event"`
	API   API   `json:"api"`
}

type Step func(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error)

type InternalState struct {
}

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
}

// func StepLogin

// func RunApplication() {
// 	authService := NewAuthService()
// 	app := Application{
// 		API: API{
// 			Endpoints: []Endpoint{
// 				{
// 					URL:    "/user",
// 					Method: []string{"POST"},
// 					Steps: []Step{
// 						authService.Login,
// 					},
// 				},
// 			},
// 		},
// 	}

// 	app.Event.User.UserCreated.Email = "arief@gmail.com"

// 	js, _ := json.MarshalIndent(app, "", "  ")
// 	log.Println(string(js))
// }
