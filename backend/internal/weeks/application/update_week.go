package application

import (
	"time"

	period_domain "backend/internal/periods/domain"
	week_domain "backend/internal/weeks/domain"
)

type UpdateWeekInput struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type UpdateWeekOutput struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type UpdateWeek struct {
	weekRepository   week_domain.WeekRepository
	periodRepository period_domain.PeriodRepository
}

func NewUpdateWeek(weekRepo week_domain.WeekRepository, periodRepo period_domain.PeriodRepository) *UpdateWeek {
	return &UpdateWeek{
		weekRepository:   weekRepo,
		periodRepository: periodRepo,
	}
}

func (uc *UpdateWeek) Execute(input UpdateWeekInput) (*UpdateWeekOutput, error) {
	week, err := uc.weekRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	_, err = uc.periodRepository.FindByID(input.PeriodID)
	if err != nil {
		return nil, err
	}

	err = week.UpdateWeek(input.Number, input.InitialDate, input.FinalDate, input.PeriodID)
	if err != nil {
		return nil, err
	}

	if err := uc.weekRepository.Update(week); err != nil {
		return nil, err
	}

	return &UpdateWeekOutput{
		ID:          week.ID,
		Number:      week.Number,
		InitialDate: week.InitialDate,
		FinalDate:   week.FinalDate,
		PeriodID:    week.PeriodID,
	}, nil
}
