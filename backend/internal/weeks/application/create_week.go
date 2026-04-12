package application

import (
	"time"

	period_domain "backend/internal/periods/domain"
	week_domain "backend/internal/weeks/domain"
)

type CreateWeekInput struct {
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type CreateWeekOutput struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type CreateWeek struct {
	weekRepository   week_domain.WeekRepository
	periodRepository period_domain.PeriodRepository
}

func NewCreateWeek(weekRepo week_domain.WeekRepository, periodRepo period_domain.PeriodRepository) *CreateWeek {
	return &CreateWeek{
		weekRepository:   weekRepo,
		periodRepository: periodRepo,
	}
}

func (uc *CreateWeek) Execute(input CreateWeekInput) (*CreateWeekOutput, error) {
	_, err := uc.periodRepository.FindByID(input.PeriodID)
	if err != nil {
		return nil, err
	}

	week, err := week_domain.NewWeek(
		input.Number,
		input.InitialDate,
		input.FinalDate,
		input.PeriodID,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.weekRepository.Create(week); err != nil {
		return nil, err
	}

	return &CreateWeekOutput{
		ID:          week.ID,
		Number:      week.Number,
		InitialDate: week.InitialDate,
		FinalDate:   week.FinalDate,
		PeriodID:    week.PeriodID,
	}, nil
}
