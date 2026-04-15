package application

import (
	"backend/internal/tasks/domain"
	"time"
)

type UpdateTaskInput struct {
	ID           uint
	AssignmentID uint
	WeekID       uint
	Title        string
	Description  string
	Status       domain.TaskStatus
	SpentHours   int
	Observations string
	Attachments  []TaskAttachmentInput
}

type UpdateTask struct {
	repository           domain.TaskRepository
	assignmentRepository TaskAssignmentRepository
	workspaceRepository  TaskWorkspaceRepository
	weekRepository       TaskWeekRepository
	now                  func() time.Time
}

func NewUpdateTask(
	repo domain.TaskRepository,
	assignmentRepo TaskAssignmentRepository,
	workspaceRepo TaskWorkspaceRepository,
	weekRepo TaskWeekRepository,
	now func() time.Time,
) *UpdateTask {
	if now == nil {
		now = time.Now
	}

	return &UpdateTask{
		repository:           repo,
		assignmentRepository: assignmentRepo,
		workspaceRepository:  workspaceRepo,
		weekRepository:       weekRepo,
		now:                  now,
	}
}

func (uc *UpdateTask) Execute(input UpdateTaskInput) (*TaskOutput, error) {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if task.AssignmentID != input.AssignmentID {
		return nil, domain.ErrTaskAssignmentChangeForbidden
	}
	if task.WeekID != input.WeekID {
		return nil, domain.ErrTaskWeekChangeForbidden
	}

	taskContext, err := loadTaskContext(
		uc.assignmentRepository,
		uc.workspaceRepository,
		uc.weekRepository,
		input.AssignmentID,
		input.WeekID,
	)
	if err != nil {
		return nil, err
	}
	if task.Late {
		return nil, domain.ErrTaskLateUpdateForbidden
	}
	if !isActiveWeek(taskContext.weekStartDate, taskContext.weekFinalDate, uc.now()) {
		return nil, domain.ErrTaskLateUpdateForbidden
	}

	attachments := make([]domain.TaskAttachment, len(input.Attachments))
	for i, attachment := range input.Attachments {
		attachments[i] = domain.TaskAttachment{Path: attachment.Path}
	}

	if err := task.UpdateTask(
		taskContext.assignment.ID,
		taskContext.week.ID,
		input.Title,
		input.Description,
		input.Status,
		input.SpentHours,
		input.Observations,
		attachments,
		taskContext.weekStartDate,
		task.Late,
	); err != nil {
		return nil, err
	}
	task.UserID = taskContext.assignment.UserID

	if err := uc.repository.Update(task); err != nil {
		return nil, err
	}

	return newTaskOutput(task), nil
}
