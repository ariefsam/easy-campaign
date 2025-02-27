package campaign

import (
	"campaign/campaign"
	"campaign/dto"
	"campaign/eventstore"
	"campaign/logger"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CampaignService struct {
	eventStore     eventStore
	campaignReader campaignReader
	campaignWriter campaignWriter
}

type campaignReader interface {
	GetCampaign(ctx context.Context, id uint) (campaign *campaign.Campaign, err error)
	GetAllCampaigns(ctx context.Context) (campaigns []campaign.Campaign, err error)
}

type campaignWriter interface {
	Project(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) (err error)
}

func NewCampaignService() (a *CampaignService, err error) {
	ctx := context.TODO()

	eventStore, _ := eventstore.New(ctx)
	campaignService, _ := campaign.New()
	return &CampaignService{
		eventStore:     eventStore,
		campaignReader: campaignService,
		campaignWriter: campaignService,
	}, nil
}

func (campaign *CampaignService) SetEventStore(eventStore eventStore) {
	campaign.eventStore = eventStore
}

func (campaign *CampaignService) CreateCampaign(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	event := dto.Event{}
	event.Campaign.CampaignCreated.UserID = payload.CreateCampaignRequest.UserID
	event.Campaign.CampaignCreated.Name = payload.CreateCampaignRequest.Name
	event.Campaign.CampaignCreated.Description = payload.CreateCampaignRequest.Description
	event.Campaign.CampaignCreated.StartDate = payload.CreateCampaignRequest.StartDate
	event.Campaign.CampaignCreated.EndDate = payload.CreateCampaignRequest.EndDate
	event.Campaign.CampaignCreated.Budget = payload.CreateCampaignRequest.Budget

	err = campaign.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Campaign.Name = payload.CreateCampaignRequest.Name
	resp.Campaign.Description = payload.CreateCampaignRequest.Description
	resp.Campaign.StartDate = payload.CreateCampaignRequest.StartDate
	resp.Campaign.EndDate = payload.CreateCampaignRequest.EndDate
	resp.Campaign.Budget = payload.CreateCampaignRequest.Budget
	resp.StatusCode = http.StatusCreated
	return
}

func (campaign *CampaignService) GetAllCampaigns(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	campaigns, err := campaign.campaignReader.GetAllCampaigns(ctx)
	if err != nil {
		logger.Error(err)
		return
	}

	var campaignResponses []struct {
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
	}

	for _, camp := range campaigns {
		campaignResponses = append(campaignResponses, struct {
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
		}{
			ID:              camp.ID,
			UserID:          camp.UserID,
			Name:            camp.Name,
			Description:     camp.Description,
			StartDate:       camp.StartDate,
			EndDate:         camp.EndDate,
			Budget:          camp.Budget,
			Status:          camp.Status,
			ChangeStartDate: camp.ChangeStartDate,
			ChangeEndDate:   camp.ChangeEndDate,
			ChangeBudget:    camp.ChangeBudget,
			CreatedAt:       camp.CreatedAt,
			UpdatedAt:       camp.UpdatedAt,
		})
	}

	resp.Campaigns = campaignResponses
	resp.StatusCode = http.StatusOK

	return
}

func (campaign *CampaignService) UpdateCampaign(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	event := dto.Event{}
	event.Campaign.CampaignUpdated.ID = payload.UpdateCampaignRequest.ID
	event.Campaign.CampaignUpdated.Name = payload.UpdateCampaignRequest.Name
	event.Campaign.CampaignUpdated.Description = payload.UpdateCampaignRequest.Description
	event.Campaign.CampaignUpdated.StartDate = payload.UpdateCampaignRequest.StartDate
	event.Campaign.CampaignUpdated.EndDate = payload.UpdateCampaignRequest.EndDate
	event.Campaign.CampaignUpdated.Budget = payload.UpdateCampaignRequest.Budget
	event.Campaign.CampaignUpdated.Status = payload.UpdateCampaignRequest.Status
	event.Campaign.CampaignUpdated.ChangeStartDate = payload.UpdateCampaignRequest.ChangeStartDate
	event.Campaign.CampaignUpdated.ChangeEndDate = payload.UpdateCampaignRequest.ChangeEndDate
	event.Campaign.CampaignUpdated.ChangeBudget = payload.UpdateCampaignRequest.ChangeBudget

	err = campaign.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}
	resp.Campaign.Name = payload.UpdateCampaignRequest.Name
	resp.Campaign.Description = payload.UpdateCampaignRequest.Description
	resp.Campaign.StartDate = payload.UpdateCampaignRequest.StartDate
	resp.Campaign.EndDate = payload.UpdateCampaignRequest.EndDate
	resp.Campaign.Budget = payload.UpdateCampaignRequest.Budget
	resp.Campaign.Status = payload.UpdateCampaignRequest.Status
	resp.Campaign.ChangeStartDate = payload.UpdateCampaignRequest.ChangeStartDate
	resp.Campaign.ChangeEndDate = payload.UpdateCampaignRequest.ChangeEndDate
	resp.Campaign.ChangeBudget = payload.UpdateCampaignRequest.ChangeBudget
	resp.StatusCode = http.StatusOK
	return
}

func (campaign *CampaignService) GetCampaign(ctx context.Context, id uint, resp *Response) (err error) {
	campaignData, err := campaign.campaignReader.GetCampaign(ctx, id)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Campaign.ID = campaignData.ID
	resp.Campaign.UserID = campaignData.UserID
	resp.Campaign.Name = campaignData.Name
	resp.Campaign.Description = campaignData.Description
	resp.Campaign.StartDate = campaignData.StartDate
	resp.Campaign.EndDate = campaignData.EndDate
	resp.Campaign.Budget = campaignData.Budget
	resp.Campaign.Status = campaignData.Status
	resp.Campaign.ChangeStartDate = campaignData.ChangeStartDate
	resp.Campaign.ChangeEndDate = campaignData.ChangeEndDate
	resp.Campaign.ChangeBudget = campaignData.ChangeBudget
	resp.Campaign.CreatedAt = campaignData.CreatedAt
	resp.Campaign.UpdatedAt = campaignData.UpdatedAt
	return
}

func (campaign *CampaignService) DeleteCampaign(ctx context.Context, id uint) (err error) {
	event := dto.Event{}
	event.Campaign.CampaignDeleted.ID = id

	err = campaign.campaignWriter.Project(ctx, "", event, time.Now())
	if err != nil {
		logger.Error(err)
		return
	}
	err = campaign.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}
	return
}

type CampaignHandler struct {
	Service *CampaignService
}

func NewCampaignHandler(service *CampaignService) CampaignHandler {
	return CampaignHandler{Service: service}
}

// GetCampaign Handler
func (h *CampaignHandler) GetCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid campaign ID", http.StatusBadRequest)
		return
	}

	var resp Response
	err = h.Service.GetCampaign(ctx, uint(id), &resp)
	if err != nil {
		http.Error(w, "Campaign not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *CampaignHandler) DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid campaign ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteCampaign(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
