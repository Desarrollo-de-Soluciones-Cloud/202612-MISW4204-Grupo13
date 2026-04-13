package application

import (
	"backend/internal/periods/domain"
)

type ClosePeriodInput struct {
	ID uint
}

type ClosePeriodOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          string              `json:"initial_date"`
	FinalDate            string              `json:"final_date"`
	InscriptionFinalDate string              `json:"inscription_final_date"`
	WeeksCount           int                 `json:"weeks_count"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type ClosePeriod struct {
	repository domain.PeriodRepository
}

func NewClosePeriod(repo domain.PeriodRepository) *ClosePeriod {
	return &ClosePeriod{repository: repo}
}

func (uc *ClosePeriod) Execute(input ClosePeriodInput) (*ClosePeriodOutput, error) {
	period, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := period.ClosePeriod(); err != nil {
		return nil, err
	}

	if err := uc.repository.Update(period); err != nil {
		return nil, err
	}

	return &ClosePeriodOutput{
		ID:                   period.ID,
		Name:                 period.Name,
		InitialDate:          period.InitialDate,
		FinalDate:            period.FinalDate,
		InscriptionFinalDate: period.InscriptionFinalDate,
		WeeksCount:           period.WeeksCount,
		PeriodState:          period.PeriodState,
	}, nil
}
