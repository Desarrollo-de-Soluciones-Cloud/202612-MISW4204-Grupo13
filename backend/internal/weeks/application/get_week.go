package application

import (
	"time"

	"backend/internal/weeks/domain"
)

type GetWeekInput struct {
	ID uint `json:"id"`
}

type GetWeekOutput struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type GetWeek struct {
	weekRepository domain.WeekRepository
}

func NewGetWeek(repo domain.WeekRepository) *GetWeek {
	return &GetWeek{
		weekRepository: repo,
	}
}

func (uc *GetWeek) Execute(input GetWeekInput) (*GetWeekOutput, error) {
	week, err := uc.weekRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetWeekOutput{
		ID:          week.ID,
		Number:      week.Number,
		InitialDate: week.InitialDate,
		FinalDate:   week.FinalDate,
		PeriodID:    week.PeriodID,
	}, nil
}
