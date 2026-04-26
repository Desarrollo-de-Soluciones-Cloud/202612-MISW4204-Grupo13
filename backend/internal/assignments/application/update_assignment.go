package application

import (
	"backend/internal/assignments/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
)

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
	repository          domain.AssignmentRepository
	workspaceRepository UpdateAssignmentWorkspaceRepository
}

type UpdateAssignmentWorkspaceRepository interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

func NewUpdateAssignment(repo domain.AssignmentRepository) *UpdateAssignment {
	return &UpdateAssignment{repository: repo}
}

func (uc *UpdateAssignment) WithWorkspaceRepository(workspaceRepo UpdateAssignmentWorkspaceRepository) *UpdateAssignment {
	uc.workspaceRepository = workspaceRepo
	return uc
}

func (uc *UpdateAssignment) validateWorkspaceNotClosed(workspaceID uint) error {
	if uc.workspaceRepository != nil {
		workspace, err := uc.workspaceRepository.FindByID(workspaceID)
		if err != nil {
			if errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) {
				return domain.ErrAssignmentWorkspaceNotFound
			}
			return err
		}
		if workspace.State == workspacesDomain.ClosedState {
			return domain.ErrAssignmentWorkspaceClosed
		}
	}
	return nil
}

func (uc *UpdateAssignment) Execute(input UpdateAssignmentInput) (*UpdateAssignmentOutput, error) {
	// TODO RF05: Revisar si las validaciones RF05 deben excluir vinculaciones de espacios cerrados cuando exista esa integracion.
	// TODO RF05: Revisar aplicacion de RF05 en futuras operaciones administrativas adicionales fuera de create/update.
	assignment, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Validate workspace is not closed
	if err := uc.validateWorkspaceNotClosed(assignment.WorkspaceID); err != nil {
		return nil, err
	}

	assignments, err := uc.repository.FindAllByUserID(assignment.UserID)
	if err != nil {
		return nil, err
	}

	for _, existing := range assignments {
		if existing.ID != assignment.ID && existing.WorkspaceID == assignment.WorkspaceID && existing.Role == input.Role {
			return nil, domain.ErrAssignmentAlreadyExists
		}
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
