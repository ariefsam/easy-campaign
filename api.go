package campaign

type API struct {
	Endpoints []Endpoint `json:"endpoints,omitzero"`
}

type Endpoint struct {
	URL    string   `json:"url,omitzero"`
	Method []string `json:"method,omitzero"`
	Steps  []Step   `json:"handler,omitzero"`
}

type Request struct {
	Login struct {
		Email       string `json:"email,omitzero"`
		Password    string `json:"password,omitzero"`
		Device      string `json:"device,omitzero"`
		Application string `json:"application,omitzero"` // admin-panel, campaign-dashboard
	} `json:"login,omitzero"`
}
