package application

import (
	"backend/internal/weeks/domain"
	"time"
)

const weekDateLayout = "2006-01-02"

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
	if err := validateCreateWeeksInput(input); err != nil {
		return nil, err
	}

	exists, err := uc.repository.ExistsByPeriodID(input.PeriodID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrWeeksAlreadyExistForPeriod
	}

	initialDate, finalDate, err := parseWeekRange(input.InitialDate, input.FinalDate)
	if err != nil {
		return nil, err
	}

	weeks := make([]domain.Week, 0, input.WeeksCount)
	for i := 0; i < input.WeeksCount; i++ {
		weekInitialDate := initialDate.AddDate(0, 0, i*7)
		weekFinalDate := weekInitialDate.AddDate(0, 0, 6)

		week, err := domain.NewWeek(
			input.PeriodID,
			i+1,
			weekInitialDate.Format(weekDateLayout),
			weekFinalDate.Format(weekDateLayout),
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

	return &CreateWeeksForPeriodOutput{Weeks: toWeekOutputs(weeks)}, nil
}

func validateCreateWeeksInput(input CreateWeeksForPeriodInput) error {
	if err := domain.ValidateWeekPeriodID(input.PeriodID); err != nil {
		return err
	}
	if err := domain.ValidateWeekInitialDate(input.InitialDate); err != nil {
		return err
	}
	if err := domain.ValidateWeekInitialDateIsMonday(input.InitialDate); err != nil {
		return err
	}
	if err := domain.ValidateWeekCount(input.WeeksCount); err != nil {
		return err
	}
	if err := domain.ValidateWeekFinalDate(input.FinalDate); err != nil {
		return err
	}
	if err := domain.ValidateWeekFinalDateIsSunday(input.FinalDate); err != nil {
		return err
	}

	return nil
}

func parseWeekRange(initialDateRaw, finalDateRaw string) (time.Time, time.Time, error) {
	initialDate, err := time.Parse(weekDateLayout, initialDateRaw)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrWeekInitialDateWrongFormat
	}

	finalDate, err := time.Parse(weekDateLayout, finalDateRaw)
	if err != nil {
		return time.Time{}, time.Time{}, domain.ErrWeekFinalDateWrongFormat
	}

	return initialDate, finalDate, nil
}

func toWeekOutputs(weeks []domain.Week) []WeekOutput {
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

	return outputWeeks
}
