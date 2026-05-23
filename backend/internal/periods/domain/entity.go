package domain

import (
	"time"
)

type Period struct {
	ID                   uint        `gorm:"primaryKey" json:"id"`
	Name                 string      `gorm:"size:100;not null" json:"name"`
	InitialDate          string      `gorm:"size:100;not null" json:"initial_date"`
	FinalDate            string      `gorm:"size:100;not null" json:"final_date"`
	InscriptionFinalDate string      `gorm:"size:100;not null" json:"inscription_final_date"`
	WeeksCount           int         `gorm:"not null" json:"weeks_count"`
	PeriodState          PeriodState `gorm:"size:20;not null" json:"period_state"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

func NewPeriod(name, initialDate string, weeksCount int, state PeriodState) (*Period, error) {
	normalizedName := NormalizePeriodName(name)

	if err := ValidatePeriodName(normalizedName); err != nil {
		return nil, err
	}
	if err := ValidatePeriodInitialDate(initialDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodInitialDateIsMonday(initialDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodInitialDateIsFuture(initialDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodWeeksCount(weeksCount); err != nil {
		return nil, err
	}
	if err := ValidatePeriodState(state); err != nil {
		return nil, err
	}

	// Calculate final date: initial_date + (weeks_count * 7) - 1 day
	finalDate, err := CalculatePeriodFinalDate(initialDate, weeksCount)
	if err != nil {
		return nil, err
	}

	// Calculate inscription final date: initial_date - 1 day
	inscriptionFinalDate, err := CalculatePeriodInscriptionFinalDate(initialDate)
	if err != nil {
		return nil, err
	}

	return &Period{
		Name:                 normalizedName,
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionFinalDate,
		WeeksCount:           weeksCount,
		PeriodState:          state,
	}, nil
}

// UpdatePeriodName updates only the name of an existing period
func (p *Period) UpdatePeriodName(name string) error {
	normalizedName := NormalizePeriodName(name)

	if err := ValidatePeriodName(normalizedName); err != nil {
		return err
	}

	p.Name = normalizedName
	return nil
}

// ClosePeriod changes the period state from Active to Closed
func (p *Period) ClosePeriod() error {
	if p.PeriodState != ActivePeriod {
		return ErrPeriodStateInvalid
	}
	p.PeriodState = ClosedPeriod
	return nil
}

// UpdatePeriod updates period with full details (legacy method - kept for compatibility)
func (p *Period) UpdatePeriod(name, initialDate string, weeksCount int, state PeriodState) error {
	// Validate that period hasn't started yet
	parsedInitialDate, err := time.Parse(periodDateLayout, initialDate)
	if err != nil {
		return ErrPeriodInitialDateWrongFormat
	}
	now := time.Now()
	if parsedInitialDate.Before(now) && parsedInitialDate.Format(periodDateLayout) != now.Format(periodDateLayout) {
		return ErrPeriodCannotBeUpdated
	}

	normalizedName := NormalizePeriodName(name)

	if err := ValidatePeriodName(normalizedName); err != nil {
		return err
	}
	if err := ValidatePeriodInitialDate(initialDate); err != nil {
		return err
	}
	if err := ValidatePeriodInitialDateIsMonday(initialDate); err != nil {
		return err
	}
	if err := ValidatePeriodWeeksCount(weeksCount); err != nil {
		return err
	}
	if err := ValidatePeriodState(state); err != nil {
		return err
	}

	// Calculate final date
	finalDate, err := CalculatePeriodFinalDate(initialDate, weeksCount)
	if err != nil {
		return err
	}

	// Calculate inscription final date
	inscriptionFinalDate, err := CalculatePeriodInscriptionFinalDate(initialDate)
	if err != nil {
		return err
	}

	p.Name = normalizedName
	p.InitialDate = initialDate
	p.FinalDate = finalDate
	p.InscriptionFinalDate = inscriptionFinalDate
	p.WeeksCount = weeksCount
	p.PeriodState = state

	return nil
}
