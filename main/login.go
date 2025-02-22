package main

import (
	"campaign"
	"campaign/apperror"
	"campaign/logger"
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func (h *campaignHandler) stepHandler(steps []campaign.Step) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := campaign.Request{}
		json.NewDecoder(r.Body).Decode(&payload)

		queryParams := r.URL.Query()
		payload.QueryParams = make(map[string]string)
		for key, values := range queryParams {
			if len(values) > 0 {
				payload.QueryParams[key] = values[0]
			}
		}

		vars := mux.Vars(r)
		inputID := vars["id"]

		if inputID != "" {
			payload.QueryParams["id"] = inputID
		}

		cookieToken, err := r.Cookie("token")
		state := &campaign.InternalState{}
		if err == nil {
			session, err := h.authService.ParseToken(r.Context(), cookieToken.Value)
			if err != nil {
				logger.Println(err)
			}
			if session != nil {
				state.Session.UserID = session.Subject
			}
		}

		resp, err := runStep(r.Context(), &payload, state, steps)
		if err != nil {
			apperror.HandleError(w, err)
			// jsonError(w, err, http.StatusInternalServerError)
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

func runStep(ctx context.Context, payload *campaign.Request, state *campaign.InternalState, step []campaign.Step) (resp campaign.Response, err error) {

	for _, s := range step {
		err = s(ctx, payload, state, &resp)
		if err != nil {
			logger.Println(err)
			return
		}
	}

	return
}
