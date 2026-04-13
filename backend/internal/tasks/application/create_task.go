package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/tasks/domain"
	"time"
)

type CreateTaskInput struct {
	AssignmentID  uint
	WeekID        *uint
	Title         string
	Description   string
	Status        domain.TaskStatus
	SpentHours    int
	Observations  string
	WeekStartDate time.Time
}

type CreateTask struct {
	repository           domain.TaskRepository
	assignmentRepository TaskAssignmentRepository
	now                  func() time.Time
}

func NewCreateTask(repo domain.TaskRepository, assignmentRepo TaskAssignmentRepository, now func() time.Time) *CreateTask {
	if now == nil {
		now = time.Now
	}

	return &CreateTask{
		repository:           repo,
		assignmentRepository: assignmentRepo,
		now:                  now,
	}
}

func (uc *CreateTask) Execute(input CreateTaskInput) (*TaskOutput, error) {
	assignment, err := uc.assignmentRepository.FindByID(input.AssignmentID)
	if err != nil {
		if err == assignmentsDomain.ErrAssignmentNotFound {
			return nil, domain.ErrTaskAssignmentNotFound
		}
		return nil, err
	}

	normalizedWeekStartDate := domain.NormalizeWeekStartDate(input.WeekStartDate)
	late := domain.IsWeekClosed(normalizedWeekStartDate, uc.now())

	task, err := domain.NewTask(
		assignment.UserID,
		assignment.ID,
		input.WeekID,
		input.Title,
		input.Description,
		input.Status,
		input.SpentHours,
		input.Observations,
		normalizedWeekStartDate,
		late,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Create(task); err != nil {
		return nil, err
	}

	return newTaskOutput(task), nil
}
