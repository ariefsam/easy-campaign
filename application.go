package campaign

import "context"

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
