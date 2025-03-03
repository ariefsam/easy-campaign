package main

import (
	"html/template"
	"net/http"
)

func renderView(w http.ResponseWriter, r *http.Request, master string, view string, data interface{}) {
	master = "view/" + master + ".html"
	view = "view/" + view + ".html"

	tmpl, err := template.ParseFiles(master, view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
