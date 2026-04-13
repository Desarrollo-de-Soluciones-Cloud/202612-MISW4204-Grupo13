package application

import "backend/internal/weeks/domain"

type GetWeekByPeriodAndNumberInput struct {
	PeriodID uint `json:"period_id"`
	Number   int  `json:"number"`
}

type GetWeekByPeriodAndNumber struct {
	repository domain.WeekRepository
}

func NewGetWeekByPeriodAndNumber(repo domain.WeekRepository) *GetWeekByPeriodAndNumber {
	return &GetWeekByPeriodAndNumber{repository: repo}
}

func (uc *GetWeekByPeriodAndNumber) Execute(input GetWeekByPeriodAndNumberInput) (*WeekOutput, error) {
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

	return &WeekOutput{
		ID:          week.ID,
		PeriodID:    week.PeriodID,
		Number:      week.Number,
		InitialDate: week.InitialDate,
		FinalDate:   week.FinalDate,
	}, nil
}
