package delivery

type CreateAssignmentRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	WorkspaceID uint   `json:"workspace_id" binding:"required"`
	Role        string `json:"role" binding:"required"`
	WeeklyHours int    `json:"weekly_hours" binding:"required,min=1"`
}

type AssignmentResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	WorkspaceID uint   `json:"workspace_id"`
	Role        string `json:"role"`
	WeeklyHours int    `json:"weekly_hours"`
}

type ListAssignmentsResponse struct {
	Assignments []AssignmentResponse `json:"assignments"`
}

type UpdateAssignmentRequest struct {
	Role        string `json:"role" binding:"required"`
	WeeklyHours int    `json:"weekly_hours" binding:"required,min=1"`
}
