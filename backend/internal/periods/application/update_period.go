package application

import (
	"backend/internal/periods/domain"
)

type UpdatePeriodInput struct {
	ID   uint
	Name string
}

type UpdatePeriodOutput struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          string              `json:"initial_date"`
	FinalDate            string              `json:"final_date"`
	InscriptionFinalDate string              `json:"inscription_final_date"`
	WeeksCount           int                 `json:"weeks_count"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type UpdatePeriod struct {
	repository domain.PeriodRepository
}

func NewUpdatePeriod(repo domain.PeriodRepository) *UpdatePeriod {
	return &UpdatePeriod{repository: repo}
}

func (uc *UpdatePeriod) Execute(input UpdatePeriodInput) (*UpdatePeriodOutput, error) {
	period, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := period.UpdatePeriodName(input.Name); err != nil {
		return nil, err
	}

	// Check if new name already exists
	p, err := uc.repository.FindByName(period.Name)
	if err == nil && p != nil && p.ID != period.ID {
		return nil, domain.ErrPeriodNameAlreadyExists
	}

	if err := uc.repository.Update(period); err != nil {
		return nil, err
	}

	return &UpdatePeriodOutput{
		ID:                   period.ID,
		Name:                 period.Name,
		InitialDate:          period.InitialDate,
		FinalDate:            period.FinalDate,
		InscriptionFinalDate: period.InscriptionFinalDate,
		WeeksCount:           period.WeeksCount,
		PeriodState:          period.PeriodState,
	}, nil
}
