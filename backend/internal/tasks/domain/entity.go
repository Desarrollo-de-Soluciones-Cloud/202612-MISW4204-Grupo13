package domain

import (
	"fmt"
	"time"
)

type TaskAttachment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FilePath    string `json:"file_path"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type Task struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	AssignmentID  uint       `gorm:"not null;index" json:"assignment_id"`
	WeekID        *uint      `gorm:"default:null" json:"week_id"`
	Title         string     `gorm:"size:255;not null" json:"title"`
	Description   string     `gorm:"type:text;not null" json:"description"`
	Status        TaskStatus `gorm:"size:30;not null" json:"status"`
	SpentHours    int        `gorm:"not null" json:"spent_hours"`
	Observations  string     `gorm:"type:text;not null" json:"observations"`
	WeekStartDate time.Time  `gorm:"type:date;not null" json:"week_start_date"`
	Late          bool       `gorm:"not null;default:false" json:"late"`
	Attachments   []TaskAttachment `gorm:"serializer:json;type:text" json:"attachments"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type TaskInput struct {
	UserID        uint
	AssignmentID  uint
	WeekID        *uint
	Title         string
	Description   string
	Status        TaskStatus
	SpentHours    int
	Observations  string
	WeekStartDate time.Time
	Late          bool
	Attachments   []TaskAttachment
}

func NewTask(input TaskInput) (*Task, error) {
	normalized := normalizeTaskInput(input)

	if err := ValidateTaskAssignmentID(normalized.AssignmentID); err != nil {
		return nil, err
	}
	if err := ValidateTaskTitle(normalized.Title); err != nil {
		return nil, err
	}
	if err := ValidateTaskDescription(normalized.Description); err != nil {
		return nil, err
	}
	if err := ValidateTaskStatus(normalized.Status); err != nil {
		return nil, err
	}
	if err := ValidateTaskSpentHours(normalized.SpentHours); err != nil {
		return nil, err
	}
	if err := ValidateTaskWeekStartDate(normalized.WeekStartDate); err != nil {
		return nil, err
	}

	return &Task{
		UserID:        normalized.UserID,
		AssignmentID:  normalized.AssignmentID,
		WeekID:        normalized.WeekID,
		Title:         normalized.Title,
		Description:   normalized.Description,
		Status:        normalized.Status,
		SpentHours:    normalized.SpentHours,
		Observations:  normalized.Observations,
		WeekStartDate: normalized.WeekStartDate,
		Late:          normalized.Late,
		Attachments:   normalizeTaskAttachments(normalized.Attachments),
	}, nil
}

func (t *Task) UpdateTask(input TaskInput) error {
	if t.Late {
		return ErrTaskLateUpdateForbidden
	}

	normalized := normalizeTaskInput(input)

	if err := ValidateTaskAssignmentID(normalized.AssignmentID); err != nil {
		return err
	}
	if err := ValidateTaskTitle(normalized.Title); err != nil {
		return err
	}
	if err := ValidateTaskDescription(normalized.Description); err != nil {
		return err
	}
	if err := ValidateTaskStatus(normalized.Status); err != nil {
		return err
	}
	if err := ValidateTaskSpentHours(normalized.SpentHours); err != nil {
		return err
	}
	if err := ValidateTaskWeekStartDate(normalized.WeekStartDate); err != nil {
		return err
	}

	t.AssignmentID = normalized.AssignmentID
	t.Title = normalized.Title
	t.Description = normalized.Description
	t.Status = normalized.Status
	t.SpentHours = normalized.SpentHours
	t.Observations = normalized.Observations
	t.WeekStartDate = normalized.WeekStartDate
	t.Late = normalized.Late
	t.Attachments = normalizeTaskAttachments(normalized.Attachments)

	return nil
}

func normalizeTaskInput(input TaskInput) TaskInput {
	input.Title = NormalizeTaskTitle(input.Title)
	input.Description = NormalizeTaskDescription(input.Description)
	input.Status = NormalizeTaskStatus(input.Status)
	input.Observations = NormalizeTaskObservations(input.Observations)
	input.WeekStartDate = NormalizeWeekStartDate(input.WeekStartDate)
	return input
}

func (t *Task) CanDelete(now time.Time) error {
	if !IsWeekActive(t.WeekStartDate, now) {
		return ErrTaskDeleteForbidden
	}
	return nil
}

func normalizeTaskAttachments(attachments []TaskAttachment) []TaskAttachment {
	if len(attachments) == 0 {
		return []TaskAttachment{}
	}

	result := make([]TaskAttachment, 0, len(attachments))
	for _, attachment := range attachments {
		attachmentID := attachment.ID
		if attachmentID == "" {
			attachmentID = fmt.Sprintf("legacy_%s", attachment.FilePath)
		}

		result = append(result, TaskAttachment{
			ID:          attachmentID,
			Name:        attachment.Name,
			FilePath:    attachment.FilePath,
			ContentType: attachment.ContentType,
			Size:        attachment.Size,
		})
	}

	return result
}

func NormalizeTaskAttachmentsForPersistence(attachments []TaskAttachment) []TaskAttachment {
	return normalizeTaskAttachments(attachments)
}
