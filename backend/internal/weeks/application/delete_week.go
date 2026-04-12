package application

import (
	"backend/internal/weeks/domain"
)

type DeleteWeekInput struct {
	ID uint `json:"id"`
}

type DeleteWeek struct {
	repository domain.WeekRepository
}

func NewDeleteWeek(repo domain.WeekRepository) *DeleteWeek {
	return &DeleteWeek{
		repository: repo,
	}
}

func (uc *DeleteWeek) Execute(input DeleteWeekInput) error {
	_, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return err
	}

	return uc.repository.Delete(input.ID)
}
