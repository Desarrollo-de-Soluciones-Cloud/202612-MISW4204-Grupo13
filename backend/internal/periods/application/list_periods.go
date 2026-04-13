package application

import (
	"backend/internal/periods/domain"
)

type ListPeriodsOutput struct {
	Periods []PeriodDTO `json:"periods"`
}

type PeriodDTO struct {
	ID                   uint                `json:"id"`
	Name                 string              `json:"name"`
	InitialDate          string              `json:"initial_date"`
	FinalDate            string              `json:"final_date"`
	InscriptionFinalDate string              `json:"inscription_final_date"`
	PeriodState          domain.PeriodState  `json:"period_state"`
}

type ListPeriods struct {
	repository domain.PeriodRepository
}

func NewListPeriods(repo domain.PeriodRepository) *ListPeriods {
	return &ListPeriods{repository: repo}
}

func (uc *ListPeriods) Execute() (*ListPeriodsOutput, error) {
	periods, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]PeriodDTO, len(periods))
	for i, p := range periods {
		result[i] = PeriodDTO{
			ID:                   p.ID,
			Name:                 p.Name,
			InitialDate:          p.InitialDate,
			FinalDate:            p.FinalDate,
			InscriptionFinalDate: p.InscriptionFinalDate,
			PeriodState:          p.PeriodState,
		}
	}
	return &ListPeriodsOutput{Periods: result}, nil
}
