package application

import (
	"backend/internal/tasks/domain"
	"time"
)

type TaskAttachmentInput struct {
	Path string
}

type CreateTaskInput struct {
	AssignmentID uint
	WeekID       uint
	Title        string
	Description  string
	Status       domain.TaskStatus
	SpentHours   int
	Observations string
	Attachments  []TaskAttachmentInput
}

type CreateTask struct {
	repository           domain.TaskRepository
	assignmentRepository TaskAssignmentRepository
	workspaceRepository  TaskWorkspaceRepository
	weekRepository       TaskWeekRepository
	now                  func() time.Time
}

func NewCreateTask(
	repo domain.TaskRepository,
	assignmentRepo TaskAssignmentRepository,
	workspaceRepo TaskWorkspaceRepository,
	weekRepo TaskWeekRepository,
	now func() time.Time,
) *CreateTask {
	if now == nil {
		now = time.Now
	}

	return &CreateTask{
		repository:           repo,
		assignmentRepository: assignmentRepo,
		workspaceRepository:  workspaceRepo,
		weekRepository:       weekRepo,
		now:                  now,
	}
}

func (uc *CreateTask) Execute(input CreateTaskInput) (*TaskOutput, error) {
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

	late := isClosedWeek(taskContext.weekFinalDate, uc.now())
	attachments := make([]domain.TaskAttachment, len(input.Attachments))
	for i, attachment := range input.Attachments {
		attachments[i] = domain.TaskAttachment{Path: attachment.Path}
	}

	task, err := domain.NewTask(
		taskContext.assignment.UserID,
		taskContext.assignment.ID,
		taskContext.week.ID,
		input.Title,
		input.Description,
		input.Status,
		input.SpentHours,
		input.Observations,
		attachments,
		taskContext.weekStartDate,
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
