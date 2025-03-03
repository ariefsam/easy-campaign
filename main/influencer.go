package main

import (
	"net/http"
)

func influencerView(w http.ResponseWriter, r *http.Request) {
	renderView(w, r, "master", "influencer", nil)
}
