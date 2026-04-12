package application

import "backend/internal/periods/domain"

type DeletePeriodInput struct {
	ID uint
}

type DeletePeriod struct {
	repository domain.PeriodRepository
}

func NewDeletePeriod(repo domain.PeriodRepository) *DeletePeriod {
	return &DeletePeriod{repository: repo}
}

func (uc *DeletePeriod) Execute(input DeletePeriodInput) error {
	period, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return err
	}

	if period.PeriodState == domain.ActivePeriod {
		return domain.ErrPeriodCannotBeDeleted
	}

	return uc.repository.Delete(input.ID)
}
