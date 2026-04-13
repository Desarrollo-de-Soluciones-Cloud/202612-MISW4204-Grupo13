package delivery

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
}

type ListTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
