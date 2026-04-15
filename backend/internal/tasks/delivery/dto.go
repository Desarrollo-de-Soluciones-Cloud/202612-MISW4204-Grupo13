package delivery

type CreateTaskRequest struct {
	AssignmentID uint                    `json:"assignment_id" binding:"required"`
	WeekID       uint                    `json:"week_id" binding:"required"`
	Title        string                  `json:"title" binding:"required"`
	Description  string                  `json:"description" binding:"required"`
	Status       string                  `json:"status" binding:"required"`
	SpentHours   int                     `json:"spent_hours" binding:"required"`
	Observations string                  `json:"observations"`
	Attachments  []TaskAttachmentRequest `json:"attachments"`
}

type UpdateTaskRequest struct {
	AssignmentID uint                    `json:"assignment_id" binding:"required"`
	WeekID       uint                    `json:"week_id" binding:"required"`
	Title        string                  `json:"title" binding:"required"`
	Description  string                  `json:"description" binding:"required"`
	Status       string                  `json:"status" binding:"required"`
	SpentHours   int                     `json:"spent_hours" binding:"required"`
	Observations string                  `json:"observations"`
	Attachments  []TaskAttachmentRequest `json:"attachments"`
}

type TaskAttachmentRequest struct {
	Path string `json:"path"`
}

type TaskAttachmentResponse struct {
	ID   uint   `json:"id"`
	Path string `json:"path"`
}

type TaskResponse struct {
	ID            uint                     `json:"id"`
	UserID        uint                     `json:"user_id"`
	AssignmentID  uint                     `json:"assignment_id"`
	WeekID        uint                     `json:"week_id"`
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	Status        string                   `json:"status"`
	SpentHours    int                      `json:"spent_hours"`
	Observations  string                   `json:"observations"`
	WeekStartDate string                   `json:"week_start_date"`
	Late          bool                     `json:"late"`
	Attachments   []TaskAttachmentResponse `json:"attachments"`
}

type ListTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
