package campaign

import (
	"campaign/dto"
	"campaign/idgenerator"
	"campaign/logger"
	"campaign/report"
	"context"
	"errors"
)

type influencerProjection interface {
	FetchInfluencers() (influencers []report.Influencer, err error)
	GetInfluencer(influencerID string) (influencer *report.Influencer, err error)
}

type InfluencerService struct {
	eventStore    eventStore
	idGenerator   idGenerator
	reportService influencerProjection
}

func NewInfluencerService(es eventStore) *InfluencerService {
	defaultIDGenerator := idgenerator.New()
	return &InfluencerService{
		eventStore:  es,
		idGenerator: defaultIDGenerator,
	}
}

func (s *InfluencerService) SetReportService(reportService influencerProjection) {
	s.reportService = reportService
}

func (s *InfluencerService) SetIDGenerator(idGenerator idGenerator) {
	s.idGenerator = idGenerator
}

func (s *InfluencerService) CreateInfluencer(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	influencerID := s.idGenerator.Generate(ctx)
	event := dto.Event{}
	event.Influencer.InfluencerCreated.InfluencerID = influencerID
	event.Influencer.InfluencerCreated.Name = payload.CreateInfluencerRequest.Name
	event.Influencer.InfluencerCreated.InstagramUsername = payload.CreateInfluencerRequest.InstagramUsername
	event.Influencer.InfluencerCreated.TiktokUsername = payload.CreateInfluencerRequest.TiktokUsername
	event.Influencer.InfluencerCreated.CreatedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Influencer.InfluencerID = influencerID
	resp.Influencer.Name = payload.CreateInfluencerRequest.Name

	return
}

func (s *InfluencerService) FetchInfluencers(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	if s.reportService == nil {
		err = errors.New("report service is not set")
		logger.Error(err)
		return
	}

	influencers, err := s.reportService.FetchInfluencers()
	if err != nil {
		logger.Error(err)
		return
	}

	influencer := make([]Influencer, len(influencers))
	for i, v := range influencers {
		influencer[i] = Influencer{
			InfluencerID:               v.InfluencerID,
			Name:                       v.Name,
			InstagramUsername:          v.InstagramUsername,
			TiktokUsername:             v.TiktokUsername,
			IsInstagramUsernameValid:   v.IsInstagramUsernameValid,
			IsTiktokUsernameValid:      v.IsTiktokUsernameValid,
			LastCheckInstagramUsername: v.LastCheckInstagramUsername,
			LastCheckTiktokUsername:    v.LastCheckTiktokUsername,
		}
	}

	resp.Influencers = influencer

	return
}

func (s *InfluencerService) UpdateInfluencer(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {
	if s.reportService == nil {
		err = errors.New("report service is not set")
		logger.Error(err)
		return
	}
	currentInfluencer, err := s.reportService.GetInfluencer(payload.UpdateInfluencerRequest.InfluencerID)
	if err != nil {
		logger.Error(err)
		return
	}

	if currentInfluencer == nil {
		err = errors.New("influencer not found")
		logger.Error(err)
		return
	}

	event := dto.Event{}
	event.Influencer.InfluencerUpdated.InfluencerID = payload.UpdateInfluencerRequest.InfluencerID
	event.Influencer.InfluencerUpdated.Name = payload.UpdateInfluencerRequest.Name
	event.Influencer.InfluencerUpdated.UpdatedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Influencer.InfluencerID = payload.UpdateInfluencerRequest.InfluencerID
	resp.Influencer.Name = payload.UpdateInfluencerRequest.Name

	return
}

func (s *InfluencerService) DeleteInfluencer(ctx context.Context, payload *Request, state *InternalState, resp *Response) (err error) {

	currentInfluencer, err := s.reportService.GetInfluencer(payload.DeleteInfluencerRequest.InfluencerID)
	if err != nil {
		logger.Error(err)
		return
	}

	if currentInfluencer == nil {
		err = errors.New("influencer not found")
		logger.Error(err)
		return
	}

	event := dto.Event{}
	event.Influencer.InfluencerDeleted.InfluencerID = payload.DeleteInfluencerRequest.InfluencerID
	event.Influencer.InfluencerDeleted.DeletedBy = state.Session.UserID
	err = s.eventStore.Save(ctx, event)
	if err != nil {
		logger.Error(err)
		return
	}

	resp.Influencer.InfluencerID = payload.DeleteInfluencerRequest.InfluencerID

	return
}
