package application

import "backend/internal/tasks/domain"

type GetTaskByIDInput struct {
	ID uint
}

type GetTaskByID struct {
	repository domain.TaskRepository
}

func NewGetTaskByID(repo domain.TaskRepository) *GetTaskByID {
	return &GetTaskByID{repository: repo}
}

func (uc *GetTaskByID) Execute(input GetTaskByIDInput) (*TaskOutput, error) {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return newTaskOutput(task), nil
}
