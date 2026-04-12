package domain

import "time"

type Week struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Number      int       `gorm:"not null" json:"number"`
	InitialDate time.Time `json:"initial_date"`
	FinalDate   time.Time `json:"final_date"`
	PeriodID    uint      `gorm:"not null;index" json:"period_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewWeek(number int, initialDate, finalDate time.Time, periodID uint) (*Week, error) {
	if err := ValidateWeekNumber(number); err != nil {
		return nil, err
	}

	if err := ValidateWeekInitialDate(initialDate); err != nil {
		return nil, err
	}

	if err := ValidateWeekFinalDate(finalDate); err != nil {
		return nil, err
	}

	if err := ValidateWeekDateSequence(initialDate, finalDate); err != nil {
		return nil, err
	}

	if err := ValidateWeekPeriodID(periodID); err != nil {
		return nil, err
	}

	return &Week{
		Number:      number,
		InitialDate: initialDate,
		FinalDate:   finalDate,
		PeriodID:    periodID,
	}, nil
}

func (w *Week) UpdateWeek(number int, initialDate, finalDate time.Time, periodID uint) error {
	if err := ValidateWeekNumber(number); err != nil {
		return err
	}

	if err := ValidateWeekInitialDate(initialDate); err != nil {
		return err
	}

	if err := ValidateWeekFinalDate(finalDate); err != nil {
		return err
	}

	if err := ValidateWeekDateSequence(initialDate, finalDate); err != nil {
		return err
	}

	if err := ValidateWeekPeriodID(periodID); err != nil {
		return err
	}

	w.Number = number
	w.InitialDate = initialDate
	w.FinalDate = finalDate
	w.PeriodID = periodID

	return nil
}