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
	weekRepo   TaskWeekRepository
	now        func() time.Time
}

func NewDeleteTask(repo domain.TaskRepository, weekRepo TaskWeekRepository, now func() time.Time) *DeleteTask {
	if now == nil {
		now = time.Now
	}

	return &DeleteTask{
		repository: repo,
		weekRepo:   weekRepo,
		now:        now,
	}
}

func (uc *DeleteTask) Execute(input DeleteTaskInput) error {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return err
	}

	week, err := uc.weekRepo.FindByID(task.WeekID)
	if err != nil {
		if err == nil {
			return domain.ErrTaskWeekNotFound
		}
		return err
	}

	weekStartDate, err := time.Parse("2006-01-02", week.InitialDate)
	if err != nil {
		return domain.ErrTaskWeekStartDateInvalid
	}
	weekFinalDate, err := time.Parse("2006-01-02", week.FinalDate)
	if err != nil {
		return domain.ErrTaskWeekStartDateInvalid
	}

	if err := task.CanDelete(isActiveWeek(weekStartDate, weekFinalDate, uc.now())); err != nil {
		return err
	}

	return uc.repository.Delete(task.ID)
}
