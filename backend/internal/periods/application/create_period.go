package application

import (
	"backend/internal/periods/domain"
)

type CreatePeriodInput struct {
	Name        string
	InitialDate string
	WeeksCount  int
	PeriodState domain.PeriodState
}

type CreatePeriodOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          string              `json:"initial_date"`
	FinalDate            string              `json:"final_date"`
	InscriptionFinalDate string              `json:"inscription_final_date"`
	WeeksCount           int                 `json:"weeks_count"`
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
		input.WeeksCount,
		input.PeriodState,
	)
	if err != nil {
		return nil, err
	}

	p, err := uc.repository.FindByName(period.Name)
	if err == nil && p != nil {
		return nil, domain.ErrPeriodNameAlreadyExists
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
		WeeksCount:           period.WeeksCount,
		PeriodState:          period.PeriodState,
	}, nil
}