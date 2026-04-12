package application

import (
	"backend/internal/weeks/domain"
)

type ListWeeksByPeriodInput struct {
	PeriodID uint `json:"period_id"`
}

type ListWeeksByPeriodOutput struct {
	Weeks []WeekResponse `json:"weeks"`
}

type ListWeeksByPeriod struct {
	weekRepository domain.WeekRepository
}

func NewListWeeksByPeriod(repo domain.WeekRepository) *ListWeeksByPeriod {
	return &ListWeeksByPeriod{
		weekRepository: repo,
	}
}

func (uc *ListWeeksByPeriod) Execute(input ListWeeksByPeriodInput) (*ListWeeksByPeriodOutput, error) {
	weeks, err := uc.weekRepository.FindAllByPeriodID(input.PeriodID)
	if err != nil {
		return nil, err
	}

	weekResponses := make([]WeekResponse, len(weeks))
	for i, week := range weeks {
		weekResponses[i] = WeekResponse{
			ID:          week.ID,
			Number:      week.Number,
			InitialDate: week.InitialDate,
			FinalDate:   week.FinalDate,
			PeriodID:    week.PeriodID,
		}
	}

	return &ListWeeksByPeriodOutput{
		Weeks: weekResponses,
	}, nil
}
