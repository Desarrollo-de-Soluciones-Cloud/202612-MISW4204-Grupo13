package delivery

type CreateWorkspaceRequest struct {
	PeriodID     uint   `json:"period_id" binding:"required"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=course project"`
	InitialDate  string `json:"initial_date" binding:"required"`
	FinalDate    string `json:"final_date" binding:"required"`
	Observations string `json:"observations" binding:"required"`
	State        string `json:"state" binding:"required,oneof=active closed"`
}

type CreateWorkspaceResponse struct {
	ID           uint   `json:"id"`
	PeriodID     uint   `json:"period_id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	InitialDate  string `json:"initial_date"`
	FinalDate    string `json:"final_date"`
	Observations string `json:"observations"`
	State        string `json:"state"`
}

type UpdateWorkspaceRequest struct {
	PeriodID     uint   `json:"period_id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name" binding:"required"`
	Type         string `json:"type" binding:"required,oneof=course project"`
	InitialDate  string `json:"initial_date" binding:"required"`
	FinalDate    string `json:"final_date" binding:"required"`
	Observations string `json:"observations" binding:"required"`
	State        string `json:"state" binding:"required,oneof=active closed"`
}

type ListWorkspacesResponse struct {
	Workspaces []WorkspaceResponse `json:"workspaces"`
}

type WorkspaceResponse struct {
	ID           uint   `json:"id"`
	PeriodID     uint   `json:"period_id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	InitialDate  string `json:"initial_date"`
	FinalDate    string `json:"final_date"`
	Observations string `json:"observations"`
	State        string `json:"state"`
}

type MonitorAssistantResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	WeeklyHours int    `json:"weekly_hours"`
}

type ListWorkspaceMonitorsAndAssistantsResponse struct {
	Monitors   []MonitorAssistantResponse `json:"monitors"`
	Assistants []MonitorAssistantResponse `json:"assistants"`
}

type CloseWorkspaceResponse struct {
	ID           uint   `json:"id"`
	PeriodID     uint   `json:"period_id"`
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	InitialDate  string `json:"initial_date"`
	FinalDate    string `json:"final_date"`
	Observations string `json:"observations"`
	State        string `json:"state"`
}
