package application

import (
	"backend/internal/periods/domain"
	"time"
)

type CreatePeriodInput struct {
	Name                 string
	InitialDate          time.Time
	FinalDate            time.Time
	InscriptionFinalDate time.Time
	PeriodState          domain.PeriodState
}

type CreatePeriodOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          time.Time           `json:"initial_date"`
	FinalDate            time.Time           `json:"final_date"`
	InscriptionFinalDate time.Time           `json:"inscription_final_date"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type CreatePeriod struct {
	repository domain.PeriodRepository
}

func NewCreatePeriod(repo domain.PeriodRepository) *CreatePeriod {
	return &CreatePeriod{repository: repo}
}

func (uc *CreatePeriod) Execute(input CreatePeriodInput) (*CreatePeriodOutput, error) {
	period, err := domain.NewPeriod(
		input.Name,
		input.InitialDate,
		input.FinalDate,
		input.InscriptionFinalDate,
		input.PeriodState,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Create(period); err != nil {
		return nil, err
	}

	return &CreatePeriodOutput{
		ID:                   period.ID,
		Name:                 period.Name,
		InitialDate:          period.InitialDate,
		FinalDate:            period.FinalDate,
		InscriptionFinalDate: period.InscriptionFinalDate,
		PeriodState:          period.PeriodState,
	}, nil
}