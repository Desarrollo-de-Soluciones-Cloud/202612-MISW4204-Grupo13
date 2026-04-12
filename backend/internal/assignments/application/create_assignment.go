package application

import "backend/internal/assignments/domain"

type CreateAssignmentInput struct {
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type CreateAssignmentOutput struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type CreateAssignment struct {
	repository domain.AssignmentRepository
}

func NewCreateAssignment(repo domain.AssignmentRepository) *CreateAssignment {
	return &CreateAssignment{repository: repo}
}

func (uc *CreateAssignment) Execute(input CreateAssignmentInput) (*CreateAssignmentOutput, error) {
	//nolint:godox // TODO RF04: Validar que user_id exista realmente cuando se defina la integracion con users.
	//nolint:godox // TODO RF04: Validar que workspace_id exista realmente cuando el modulo de workspaces este terminado.
	assistantHours, err := uc.repository.SumWeeklyHoursByUserAndRole(input.UserID, domain.RoleAssistant)
	if err != nil {
		return nil, err
	}

	monitorHours, err := uc.repository.SumWeeklyHoursByUserAndRole(input.UserID, domain.RoleMonitor)
	if err != nil {
		return nil, err
	}

	monitorCount, err := uc.repository.CountAssignmentsByUserAndRole(input.UserID, domain.RoleMonitor)
	if err != nil {
		return nil, err
	}

	nextWorkload, err := domain.BuildWorkloadWithAssignment(domain.UserAssignmentWorkload{
		AssistantWeeklyHours: assistantHours,
		MonitorWeeklyHours:   monitorHours,
		MonitorAssignments:   monitorCount,
	}, input.Role, input.WeeklyHours)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateRF05Workload(nextWorkload); err != nil {
		return nil, err
	}

	assignment, err := domain.NewAssignment(
		input.UserID,
		input.WorkspaceID,
		input.Role,
		input.WeeklyHours,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Create(assignment); err != nil {
		return nil, err
	}

	return &CreateAssignmentOutput{
		ID:          assignment.ID,
		UserID:      assignment.UserID,
		WorkspaceID: assignment.WorkspaceID,
		Role:        assignment.Role,
		WeeklyHours: assignment.WeeklyHours,
	}, nil
}
