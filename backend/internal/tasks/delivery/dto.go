package delivery

import "backend/internal/tasks/domain"

type CreateTaskRequest struct {
	AssignmentID  uint   `json:"assignment_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Status        string `json:"status" binding:"required"`
	SpentHours    int    `json:"spent_hours" binding:"required"`
	Observations  string `json:"observations"`
	WeekStartDate string `json:"week_start_date" binding:"required"`
}

type UpdateTaskRequest struct {
	AssignmentID  uint   `json:"assignment_id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Status        string `json:"status" binding:"required"`
	SpentHours    int    `json:"spent_hours" binding:"required"`
	Observations  string `json:"observations"`
	WeekStartDate string `json:"week_start_date" binding:"required"`
}

type TaskAttachmentResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FilePath    string `json:"file_path"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type TaskResponse struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	AssignmentID  uint   `json:"assignment_id"`
	WeekID        *uint  `json:"week_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	SpentHours    int    `json:"spent_hours"`
	Observations  string `json:"observations"`
	WeekStartDate string `json:"week_start_date"`
	Late          bool   `json:"late"`
	Attachments   []TaskAttachmentResponse `json:"attachments"`
}

type ListTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

func toTaskAttachmentResponses(attachments []domain.TaskAttachment) []TaskAttachmentResponse {
	if len(attachments) == 0 {
		return []TaskAttachmentResponse{}
	}

	result := make([]TaskAttachmentResponse, 0, len(attachments))
	for _, attachment := range attachments {
		result = append(result, TaskAttachmentResponse{
			ID:          attachment.ID,
			Name:        attachment.Name,
			FilePath:    attachment.FilePath,
			ContentType: attachment.ContentType,
			Size:        attachment.Size,
		})
	}

	return result
}

