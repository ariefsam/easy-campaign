package main

import (
	"campaign"
	"campaign/logger"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func loginView(w http.ResponseWriter, r *http.Request) {
	view := "view/login.html"

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

// func (h *campaignHandler) loginPost(w http.ResponseWriter, r *http.Request) {

// 	payload := campaign.Request{}
// 	json.NewDecoder(r.Body).Decode(&payload)

// 	resp, err := runStep(r.Context(), &payload, []campaign.Step{h.authService.Login})
// 	if err != nil {
// 		jsonError(w, err, http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	err = json.NewEncoder(w).Encode(resp)
// 	if err != nil {
// 		logger.Println(err)
// 		jsonError(w, err, http.StatusInternalServerError)
// 	}
// }

func (h *campaignHandler) stepHandler(steps []campaign.Step) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := campaign.Request{}
		json.NewDecoder(r.Body).Decode(&payload)

		resp, err := runStep(r.Context(), &payload, steps)
		if err != nil {
			jsonError(w, err, http.StatusInternalServerError)
			return
		}

		if resp.Auth.Token != "" {
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   resp.Auth.Token,
				Expires: time.Now().Add(24 * time.Hour),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			logger.Println(err)
			jsonError(w, err, http.StatusInternalServerError)
		}
	}
}

func jsonError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   true,
		"message": cases.Upper(language.Indonesian).String(err.Error()),
	})
}

func runStep(ctx context.Context, payload *campaign.Request, step []campaign.Step) (resp campaign.Response, err error) {
	state := campaign.InternalState{}

	for _, s := range step {
		err = s(ctx, payload, &state, &resp)
		if err != nil {
			logger.Println(err)
			return
		}
	}

	return
}
