package application

import (
	"backend/internal/weeks/domain"
	"time"
)

type GetWeekStatusByPeriodAndNumberInput struct {
	PeriodID uint `json:"period_id"`
	Number   int  `json:"number"`
}

type GetWeekStatusByPeriodAndNumberOutput struct {
	ID          uint              `json:"id"`
	PeriodID    uint              `json:"period_id"`
	Number      int               `json:"number"`
	InitialDate string            `json:"initial_date"`
	FinalDate   string            `json:"final_date"`
	Status      domain.WeekStatus `json:"status"`
}

type GetWeekStatusByPeriodAndNumber struct {
	repository domain.WeekRepository
	now        func() time.Time
}

func NewGetWeekStatusByPeriodAndNumber(repo domain.WeekRepository) *GetWeekStatusByPeriodAndNumber {
	return &GetWeekStatusByPeriodAndNumber{
		repository: repo,
		now:        time.Now,
	}
}

func NewGetWeekStatusByPeriodAndNumberWithNow(repo domain.WeekRepository, now func() time.Time) *GetWeekStatusByPeriodAndNumber {
	return &GetWeekStatusByPeriodAndNumber{
		repository: repo,
		now:        now,
	}
}

func (uc *GetWeekStatusByPeriodAndNumber) Execute(input GetWeekStatusByPeriodAndNumberInput) (*GetWeekStatusByPeriodAndNumberOutput, error) {
	if err := domain.ValidateWeekPeriodID(input.PeriodID); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekNumber(input.Number); err != nil {
		return nil, err
	}

	week, err := uc.repository.FindByPeriodIDAndNumber(input.PeriodID, input.Number)
	if err != nil {
		return nil, err
	}

	finalDate, err := time.Parse("2006-01-02", week.FinalDate)
	if err != nil {
		return nil, domain.ErrWeekFinalDateWrongFormat
	}

	currentDate := uc.now()
	currentDate = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())

	status := domain.WeekStatusActive
	if currentDate.After(finalDate) {
		status = domain.WeekStatusClosed
	}

	return &GetWeekStatusByPeriodAndNumberOutput{
		ID:          week.ID,
		PeriodID:    week.PeriodID,
		Number:      week.Number,
		InitialDate: week.InitialDate,
		FinalDate:   week.FinalDate,
		Status:      status,
	}, nil
}
