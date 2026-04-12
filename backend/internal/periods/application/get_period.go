package application

import (
	"backend/internal/periods/domain"
)

type GetPeriodByIDInput struct {
	ID uint
}

type GetPeriodByIDOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          string              `json:"initial_date"`
	FinalDate            string              `json:"final_date"`
	InscriptionFinalDate string              `json:"inscription_final_date"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type GetPeriodByID struct {
	repository domain.PeriodRepository
}

func NewGetPeriodByID(repo domain.PeriodRepository) *GetPeriodByID {
	return &GetPeriodByID{repository: repo}
}

func (uc *GetPeriodByID) Execute(input GetPeriodByIDInput) (*GetPeriodByIDOutput, error) {
	period, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetPeriodByIDOutput{
		ID:                   period.ID,
		Name:                 period.Name,
		InitialDate:          period.InitialDate,
		FinalDate:            period.FinalDate,
		InscriptionFinalDate: period.InscriptionFinalDate,
		PeriodState:          period.PeriodState,
	}, nil
}
