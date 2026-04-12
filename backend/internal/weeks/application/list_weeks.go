package application

import (
	"time"

	"backend/internal/weeks/domain"
)

type ListWeeksOutput struct {
	Weeks []WeekResponse `json:"weeks"`
}

type WeekResponse struct {
	ID          uint      `json:"id"`
	Number      int       `json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `json:"period_id"`
}

type ListWeeks struct {
	weekRepository domain.WeekRepository
}

func NewListWeeks(repo domain.WeekRepository) *ListWeeks {
	return &ListWeeks{
		weekRepository: repo,
	}
}

func (uc *ListWeeks) Execute() (*ListWeeksOutput, error) {
	weeks, err := uc.weekRepository.FindAll()
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

	return &ListWeeksOutput{
		Weeks: weekResponses,
	}, nil
}
