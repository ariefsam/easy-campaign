package main

import (
	"html/template"
	"net/http"
)

func campaignView(w http.ResponseWriter, r *http.Request) {
	view := "view/campaign.html"

	tmpl, err := template.ParseFiles(view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
