package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/tasks/domain"
	"time"
)

type UpdateTaskInput struct {
	ID            uint
	AssignmentID  uint
	Title         string
	Description   string
	Status        domain.TaskStatus
	SpentHours    int
	Observations  string
	WeekStartDate time.Time
	Attachments   []domain.TaskAttachment
}

type UpdateTask struct {
	repository           domain.TaskRepository
	assignmentRepository TaskAssignmentRepository
	now                  func() time.Time
}

func NewUpdateTask(repo domain.TaskRepository, assignmentRepo TaskAssignmentRepository, now func() time.Time) *UpdateTask {
	if now == nil {
		now = time.Now
	}

	return &UpdateTask{
		repository:           repo,
		assignmentRepository: assignmentRepo,
		now:                  now,
	}
}

func (uc *UpdateTask) Execute(input UpdateTaskInput) (*TaskOutput, error) {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	assignment, err := uc.assignmentRepository.FindByID(input.AssignmentID)
	if err != nil {
		if err == assignmentsDomain.ErrAssignmentNotFound {
			return nil, domain.ErrTaskAssignmentNotFound
		}
		return nil, err
	}
	if domain.IsWeekClosed(task.WeekStartDate, uc.now()) {
		return nil, domain.ErrTaskLateUpdateForbidden
	}
	if task.Late {
		return nil, domain.ErrTaskLateUpdateForbidden
	}

	normalizedWeekStartDate := domain.NormalizeWeekStartDate(input.WeekStartDate)
	if domain.IsWeekClosed(normalizedWeekStartDate, uc.now()) {
		return nil, domain.ErrTaskLateUpdateForbidden
	}

	if err := task.UpdateTask(
		assignment.ID,
		input.Title,
		input.Description,
		input.Status,
		input.SpentHours,
		input.Observations,
		normalizedWeekStartDate,
		false,
		input.Attachments,
	); err != nil {
		return nil, err
	}
	task.UserID = assignment.UserID

	if err := uc.repository.Update(task); err != nil {
		return nil, err
	}

	return newTaskOutput(task), nil
}
