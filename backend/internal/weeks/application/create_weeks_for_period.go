package application

import (
	"backend/internal/weeks/domain"
	"time"
)

type CreateWeeksForPeriodInput struct {
	PeriodID    uint   `json:"period_id"`
	InitialDate string `json:"initial_date"`
	FinalDate   string `json:"final_date"`
	WeeksCount  int    `json:"weeks_count"`
}

type WeekOutput struct {
	ID          uint   `json:"id"`
	PeriodID    uint   `json:"period_id"`
	Number      int    `json:"number"`
	InitialDate string `json:"initial_date"`
	FinalDate   string `json:"final_date"`
}

type CreateWeeksForPeriodOutput struct {
	Weeks []WeekOutput `json:"weeks"`
}

type CreateWeeksForPeriod struct {
	repository domain.WeekRepository
}

func NewCreateWeeksForPeriod(repo domain.WeekRepository) *CreateWeeksForPeriod {
	return &CreateWeeksForPeriod{repository: repo}
}

func (uc *CreateWeeksForPeriod) Execute(input CreateWeeksForPeriodInput) (*CreateWeeksForPeriodOutput, error) {
	if err := domain.ValidateWeekPeriodID(input.PeriodID); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekInitialDate(input.InitialDate); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekInitialDateIsMonday(input.InitialDate); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekCount(input.WeeksCount); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekFinalDate(input.FinalDate); err != nil {
		return nil, err
	}
	if err := domain.ValidateWeekFinalDateIsSunday(input.FinalDate); err != nil {
		return nil, err
	}

	exists, err := uc.repository.ExistsByPeriodID(input.PeriodID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrWeeksAlreadyExistForPeriod
	}

	initialDate, err := time.Parse("2006-01-02", input.InitialDate)
	if err != nil {
		return nil, domain.ErrWeekInitialDateWrongFormat
	}
	finalDate, err := time.Parse("2006-01-02", input.FinalDate)
	if err != nil {
		return nil, domain.ErrWeekFinalDateWrongFormat
	}

	weeks := make([]domain.Week, 0, input.WeeksCount)
	for i := 0; i < input.WeeksCount; i++ {
		weekInitialDate := initialDate.AddDate(0, 0, i*7)
		weekFinalDate := weekInitialDate.AddDate(0, 0, 6)

		week, err := domain.NewWeek(
			input.PeriodID,
			i+1,
			weekInitialDate.Format("2006-01-02"),
			weekFinalDate.Format("2006-01-02"),
		)
		if err != nil {
			return nil, err
		}

		weeks = append(weeks, *week)
	}

	expectedFinalDate := initialDate.AddDate(0, 0, input.WeeksCount*7-1)
	if !expectedFinalDate.Equal(finalDate) {
		return nil, domain.ErrWeekFinalDateMismatch
	}

	if err := uc.repository.CreateMany(weeks); err != nil {
		return nil, err
	}

	outputWeeks := make([]WeekOutput, len(weeks))
	for i, week := range weeks {
		outputWeeks[i] = WeekOutput{
			ID:          week.ID,
			PeriodID:    week.PeriodID,
			Number:      week.Number,
			InitialDate: week.InitialDate,
			FinalDate:   week.FinalDate,
		}
	}

	return &CreateWeeksForPeriodOutput{Weeks: outputWeeks}, nil
}
