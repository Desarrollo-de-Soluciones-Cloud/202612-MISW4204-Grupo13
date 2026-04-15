package application

import (
	"backend/internal/tasks/domain"
	"time"
)

type TaskOutput struct {
	ID            uint                   `json:"id"`
	UserID        uint                   `json:"user_id"`
	AssignmentID  uint                   `json:"assignment_id"`
	WeekID        uint                   `json:"week_id"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Status        domain.TaskStatus      `json:"status"`
	SpentHours    int                    `json:"spent_hours"`
	Observations  string                 `json:"observations"`
	WeekStartDate time.Time              `json:"week_start_date"`
	Late          bool                   `json:"late"`
	Attachments   []TaskAttachmentOutput `json:"attachments"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type TaskAttachmentOutput struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

func newTaskOutput(task *domain.Task) *TaskOutput {
	attachments := make([]TaskAttachmentOutput, len(task.Attachments))
	for i, attachment := range task.Attachments {
		attachments[i] = TaskAttachmentOutput{
			ID:   attachment.ID,
			Path: attachment.Path,
		}
	}

	return &TaskOutput{
		ID:            task.ID,
		UserID:        task.UserID,
		AssignmentID:  task.AssignmentID,
		WeekID:        task.WeekID,
		Title:         task.Title,
		Description:   task.Description,
		Status:        task.Status,
		SpentHours:    task.SpentHours,
		Observations:  task.Observations,
		WeekStartDate: task.WeekStartDate,
		Late:          task.Late,
		Attachments:   attachments,
		CreatedAt:     task.CreatedAt,
		UpdatedAt:     task.UpdatedAt,
	}
}
