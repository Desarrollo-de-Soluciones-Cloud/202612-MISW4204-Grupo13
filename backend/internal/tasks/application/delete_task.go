package application

import (
	"backend/internal/tasks/domain"
	"time"
)

type DeleteTaskInput struct {
	ID uint
}

type DeleteTask struct {
	repository domain.TaskRepository
	now        func() time.Time
}

func NewDeleteTask(repo domain.TaskRepository, now func() time.Time) *DeleteTask {
	if now == nil {
		now = time.Now
	}

	return &DeleteTask{
		repository: repo,
		now:        now,
	}
}

func (uc *DeleteTask) Execute(input DeleteTaskInput) error {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return err
	}
	if err := task.CanDelete(uc.now()); err != nil {
		return err
	}

	return uc.repository.Delete(task.ID)
}
