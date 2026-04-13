package application

import "backend/internal/weeks/domain"

type ListWeeksByPeriodInput struct {
	PeriodID uint `json:"period_id"`
}

type ListWeeksByPeriodOutput struct {
	Weeks []WeekOutput `json:"weeks"`
}

type ListWeeksByPeriod struct {
	repository domain.WeekRepository
}

func NewListWeeksByPeriod(repo domain.WeekRepository) *ListWeeksByPeriod {
	return &ListWeeksByPeriod{repository: repo}
}

func (uc *ListWeeksByPeriod) Execute(input ListWeeksByPeriodInput) (*ListWeeksByPeriodOutput, error) {
	if err := domain.ValidateWeekPeriodID(input.PeriodID); err != nil {
		return nil, err
	}

	weeks, err := uc.repository.FindAllByPeriodID(input.PeriodID)
	if err != nil {
		return nil, err
	}

	output := make([]WeekOutput, len(weeks))
	for i, week := range weeks {
		output[i] = WeekOutput{
			ID:          week.ID,
			PeriodID:    week.PeriodID,
			Number:      week.Number,
			InitialDate: week.InitialDate,
			FinalDate:   week.FinalDate,
		}
	}

	return &ListWeeksByPeriodOutput{Weeks: output}, nil
}
