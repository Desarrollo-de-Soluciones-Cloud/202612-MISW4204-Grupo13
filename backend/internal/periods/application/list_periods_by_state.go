package application

import "backend/internal/periods/domain"

type ListPeriodsByStateInput struct {
	PeriodState domain.PeriodState
}

type ListPeriodsByStateOutput struct {
	Periods []PeriodDTO `json:"periods"`
}

type ListPeriodsByState struct {
	repository domain.PeriodRepository
}

func NewListPeriodsByState(repo domain.PeriodRepository) *ListPeriodsByState {
	return &ListPeriodsByState{repository: repo}
}

func (uc *ListPeriodsByState) Execute(input ListPeriodsByStateInput) (*ListPeriodsByStateOutput, error) {
	if err := domain.ValidatePeriodState(input.PeriodState); err != nil {
		return nil, err
	}

	periods, err := uc.repository.FindAllByState(input.PeriodState)
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

	return &ListPeriodsByStateOutput{Periods: result}, nil
}
