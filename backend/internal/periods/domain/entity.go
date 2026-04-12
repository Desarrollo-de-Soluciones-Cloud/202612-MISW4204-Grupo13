package domain

import (
	"time"
)

type Period struct {
	ID                   uint        `gorm:"primaryKey" json:"id"`
	Name                 string      `gorm:"size:100;not null" json:"name"`
	InitialDate          time.Time   `json:"initial_date"`
	FinalDate            time.Time   `json:"final_date"`
	InscriptionFinalDate time.Time   `json:"inscription_final_date"`
	PeriodState          PeriodState `gorm:"size:20;not null" json:"period_state"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

func NewPeriod(name string, initialDate, finalDate, inscriptionFinalDate time.Time, state PeriodState) (*Period, error) {
	normalizedName := NormalizePeriodName(name)

	if err := ValidatePeriodName(normalizedName); err != nil {
		return nil, err
	}
	if err := ValidatePeriodInitialDate(initialDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodFinalDate(finalDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodInscriptionFinalDate(inscriptionFinalDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodDateSequence(initialDate, finalDate, inscriptionFinalDate); err != nil {
		return nil, err
	}
	if err := ValidatePeriodState(state); err != nil {
		return nil, err
	}

	return &Period{
		Name:                 normalizedName,
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionFinalDate,
		PeriodState:          state,
	}, nil
}

func (p *Period) UpdatePeriod(name string, initialDate, finalDate, inscriptionFinalDate time.Time, state PeriodState) error {
	normalizedName := NormalizePeriodName(name)

	if err := ValidatePeriodName(normalizedName); err != nil {
		return err
	}
	if err := ValidatePeriodInitialDate(initialDate); err != nil {
		return err
	}
	if err := ValidatePeriodFinalDate(finalDate); err != nil {
		return err
	}
	if err := ValidatePeriodInscriptionFinalDate(inscriptionFinalDate); err != nil {
		return err
	}
	if err := ValidatePeriodDateSequence(initialDate, finalDate, inscriptionFinalDate); err != nil {
		return err
	}
	if err := ValidatePeriodState(state); err != nil {
		return err
	}

	p.Name = normalizedName
	p.InitialDate = initialDate
	p.FinalDate = finalDate
	p.InscriptionFinalDate = inscriptionFinalDate
	p.PeriodState = state

	return nil
}