package campaign

type API struct {
	Endpoints []Endpoint `json:"endpoints,omitzero"`
}

type Endpoint struct {
	URL    string   `json:"url,omitzero"`
	Method []string `json:"method,omitzero"`
	Steps  []Step   `json:"handler,omitzero"`
}
