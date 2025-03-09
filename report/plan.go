package report

import (
	"campaign/dto"
	"campaign/logger"
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	PlanID    string    `json:"plan_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedBy string    `json:"created_by"`
}

func (s *reportService) handlePlanEvent(ctx context.Context, eventID, eventName string, event dto.Event, dateTime time.Time) {
	planHandler := &planHandler{
		db: s.db,
	}
	switch eventName {
	case "plan_created":
		planHandler.handlePlanCreated(ctx, eventID, event, dateTime)
	case "plan_updated":
		planHandler.handlePlanUpdated(ctx, eventID, event, dateTime)
	case "plan_deleted":
		planHandler.handlePlanDeleted(ctx, eventID, event, dateTime)
	}
}

type planHandler struct {
	db *gorm.DB
}

func (p *planHandler) handlePlanCreated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		plan := &Plan{
			PlanID:    event.Plan.PlanCreated.PlanID,
			Name:      event.Plan.PlanCreated.Name,
			StartDate: event.Plan.PlanCreated.StartDate,
			EndDate:   event.Plan.PlanCreated.EndDate,
			CreatedBy: event.Plan.PlanCreated.CreatedBy,
		}
		if err := tx.Create(plan).Error; err != nil {
			logger.Error(err)
			return err
		}
		cursor := &Cursor{
			EventID: eventID,
		}

		if err := tx.Create(cursor).Error; err != nil {
			logger.Error(err)
		}
		return nil
	})

	if err != nil {
		logger.Error(err)
	}
}

func (p *planHandler) handlePlanUpdated(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		plan := &Plan{}
		if err := tx.Where("plan_id = ?", event.Plan.PlanUpdated.PlanID).First(plan).Error; err != nil {
			return err
		}

		plan.Name = event.Plan.PlanUpdated.Name
		plan.StartDate = event.Plan.PlanUpdated.StartDate
		plan.EndDate = event.Plan.PlanUpdated.EndDate
		plan.CreatedBy = event.Plan.PlanUpdated.UpdatedBy

		if err := tx.Save(plan).Error; err != nil {
			return err
		}

		cursor := &Cursor{
			EventID: eventID,
		}

		if err := tx.Create(cursor).Error; err != nil {
			logger.Error(err)
		}
		return nil
	})

	if err != nil {
		logger.Error(err)
	}
}

func (p *planHandler) handlePlanDeleted(ctx context.Context, eventID string, event dto.Event, dateTime time.Time) {
	err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("plan_id = ?", event.Plan.PlanDeleted.PlanID).Delete(&Plan{}).Error; err != nil {
			return err
		}

		cursor := &Cursor{
			EventID: eventID,
		}

		if err := tx.Create(cursor).Error; err != nil {
			logger.Error(err)
		}
		return nil
	})

	if err != nil {
		logger.Error(err)
	}
}

func (r *reportService) GetPlan(planID string) (*Plan, error) {
	plan := &Plan{}
	err := r.db.Where("plan_id = ?", planID).First(plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return plan, nil
}

func (r *reportService) FetchPlans() (plans []Plan, err error) {
	err = r.db.Find(&plans).Error
	if err != nil {
		err = errors.Wrap(err, "failed to fetch plans")
		logger.Println(err)
		return
	}

	return
}
