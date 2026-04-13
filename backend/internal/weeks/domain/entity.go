package domain

type Week struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	PeriodID    uint   `gorm:"not null;uniqueIndex:idx_weeks_period_number;uniqueIndex:idx_weeks_period_dates" json:"period_id"`
	Number      int    `gorm:"not null;uniqueIndex:idx_weeks_period_number" json:"number"`
	InitialDate string `gorm:"size:10;not null;uniqueIndex:idx_weeks_period_dates" json:"initial_date"`
	FinalDate   string `gorm:"size:10;not null;uniqueIndex:idx_weeks_period_dates" json:"final_date"`
}

func NewWeek(periodID uint, number int, initialDate, finalDate string) (*Week, error) {
	if err := ValidateWeekPeriodID(periodID); err != nil {
		return nil, err
	}
	if err := ValidateWeekNumber(number); err != nil {
		return nil, err
	}
	if err := ValidateWeekInitialDate(initialDate); err != nil {
		return nil, err
	}
	if err := ValidateWeekFinalDate(finalDate); err != nil {
		return nil, err
	}
	if err := ValidateWeekInitialDateIsMonday(initialDate); err != nil {
		return nil, err
	}
	if err := ValidateWeekFinalDateIsSunday(finalDate); err != nil {
		return nil, err
	}
	if err := ValidateWeekDateRange(initialDate, finalDate); err != nil {
		return nil, err
	}

	return &Week{
		PeriodID:    periodID,
		Number:      number,
		InitialDate: initialDate,
		FinalDate:   finalDate,
	}, nil
}
