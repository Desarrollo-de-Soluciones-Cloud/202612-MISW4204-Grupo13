package application

import (
	"backend/internal/periods/domain"
	"time"
)

type GetPeriodByNameInput struct {
	Name string
}

type GetPeriodByNameOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          time.Time           `json:"initial_date"`
	FinalDate            time.Time           `json:"final_date"`
	InscriptionFinalDate time.Time           `json:"inscription_final_date"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type GetPeriodByName struct {
	repository domain.PeriodRepository
}

func NewGetPeriodByName(repo domain.PeriodRepository) *GetPeriodByName {
	return &GetPeriodByName{repository: repo}
}

func (uc *GetPeriodByName) Execute(input GetPeriodByNameInput) (*GetPeriodByNameOutput, error) {
	if err := domain.ValidatePeriodName(input.Name); err != nil {
		return nil, err
	}

	period, err := uc.repository.FindByName(domain.NormalizePeriodName(input.Name))
	if err != nil {
		return nil, err
	}

	return &GetPeriodByNameOutput{
		ID:                   period.ID,
		Name:                 period.Name,
		InitialDate:          period.InitialDate,
		FinalDate:            period.FinalDate,
		InscriptionFinalDate: period.InscriptionFinalDate,
		PeriodState:          period.PeriodState,
	}, nil
}
