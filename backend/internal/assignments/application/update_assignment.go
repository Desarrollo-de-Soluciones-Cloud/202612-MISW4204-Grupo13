package application

import "backend/internal/assignments/domain"

type UpdateAssignmentInput struct {
	ID          uint
	Role        domain.AssignmentRole
	WeeklyHours int
}

type UpdateAssignmentOutput struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type UpdateAssignment struct {
	repository domain.AssignmentRepository
}

func NewUpdateAssignment(repo domain.AssignmentRepository) *UpdateAssignment {
	return &UpdateAssignment{repository: repo}
}

func (uc *UpdateAssignment) Execute(input UpdateAssignmentInput) (*UpdateAssignmentOutput, error) {
	// TODO RF05: Revisar si las validaciones RF05 deben excluir vinculaciones de espacios cerrados cuando exista esa integracion.
	// TODO RF05: Revisar aplicacion de RF05 en futuras operaciones administrativas adicionales fuera de create/update.
	assignment, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	assistantHours, err := uc.repository.SumWeeklyHoursByUserAndRole(assignment.UserID, domain.RoleAssistant)
	if err != nil {
		return nil, err
	}

	monitorHours, err := uc.repository.SumWeeklyHoursByUserAndRole(assignment.UserID, domain.RoleMonitor)
	if err != nil {
		return nil, err
	}

	monitorCount, err := uc.repository.CountAssignmentsByUserAndRole(assignment.UserID, domain.RoleMonitor)
	if err != nil {
		return nil, err
	}

	currentWorkload := domain.UserAssignmentWorkload{
		AssistantWeeklyHours: assistantHours,
		MonitorWeeklyHours:   monitorHours,
		MonitorAssignments:   monitorCount,
	}

	if assignment.Role == domain.RoleAssistant {
		currentWorkload.AssistantWeeklyHours -= assignment.WeeklyHours
	}
	if assignment.Role == domain.RoleMonitor {
		currentWorkload.MonitorWeeklyHours -= assignment.WeeklyHours
		currentWorkload.MonitorAssignments--
	}

	nextWorkload, err := domain.BuildWorkloadWithAssignment(currentWorkload, input.Role, input.WeeklyHours)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateRF05Workload(nextWorkload); err != nil {
		return nil, err
	}

	if err := assignment.UpdateAdmin(input.Role, input.WeeklyHours); err != nil {
		return nil, err
	}

	if err := uc.repository.Update(assignment); err != nil {
		return nil, err
	}

	return &UpdateAssignmentOutput{
		ID:          assignment.ID,
		UserID:      assignment.UserID,
		WorkspaceID: assignment.WorkspaceID,
		Role:        assignment.Role,
		WeeklyHours: assignment.WeeklyHours,
	}, nil
}
