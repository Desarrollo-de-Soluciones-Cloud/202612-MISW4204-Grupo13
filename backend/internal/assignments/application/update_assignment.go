package application

import (
	"backend/internal/assignments/domain"
	periodsDomain "backend/internal/periods/domain"
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
	periodRepository    UpdateAssignmentPeriodRepository
}

type UpdateAssignmentWorkspaceRepository interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type UpdateAssignmentPeriodRepository interface {
	FindByID(id uint) (*periodsDomain.Period, error)
}

func NewUpdateAssignment(repo domain.AssignmentRepository) *UpdateAssignment {
	return &UpdateAssignment{repository: repo}
}

func (uc *UpdateAssignment) WithWorkspaceRepository(workspaceRepo UpdateAssignmentWorkspaceRepository) *UpdateAssignment {
	uc.workspaceRepository = workspaceRepo
	return uc
}

func (uc *UpdateAssignment) WithPeriodRepository(periodRepo UpdateAssignmentPeriodRepository) *UpdateAssignment {
	uc.periodRepository = periodRepo
	return uc
}

func (uc *UpdateAssignment) validateWorkspaceNotClosed(workspaceID uint) error {
	workspace, err := uc.validateWorkspace(workspaceID)
	if err != nil || workspace == nil {
		return err
	}

	return uc.validateWorkspacePeriod(workspace.PeriodID)
}

func (uc *UpdateAssignment) validateWorkspace(workspaceID uint) (*workspacesDomain.Workspace, error) {
	if uc.workspaceRepository == nil {
		return nil, nil
	}

	workspace, err := uc.workspaceRepository.FindByID(workspaceID)
	if err != nil {
		if errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) {
			return nil, domain.ErrAssignmentWorkspaceNotFound
		}
		return nil, err
	}

	if workspace.State == workspacesDomain.ClosedState {
		return nil, domain.ErrAssignmentWorkspaceClosed
	}

	return workspace, nil
}

func (uc *UpdateAssignment) validateWorkspacePeriod(periodID uint) error {
	if uc.periodRepository == nil {
		return nil
	}

	period, err := uc.periodRepository.FindByID(periodID)
	if err != nil {
		if errors.Is(err, periodsDomain.ErrPeriodNotFound) {
			return domain.ErrAssignmentWorkspaceNotFound
		}
		return err
	}

	if period.PeriodState == periodsDomain.ClosedPeriod {
		return domain.ErrAssignmentPeriodClosed
	}

	return nil
}

func (uc *UpdateAssignment) ensureNoExactDuplicate(assignment *domain.Assignment, role domain.AssignmentRole) error {
	assignments, err := uc.repository.FindAllByUserID(assignment.UserID)
	if err != nil {
		return err
	}

	for _, existing := range assignments {
		if existing.ID != assignment.ID && existing.WorkspaceID == assignment.WorkspaceID && existing.Role == role {
			return domain.ErrAssignmentAlreadyExists
		}
	}

	return nil
}

func (uc *UpdateAssignment) buildCurrentWorkload(assignment *domain.Assignment) (domain.UserAssignmentWorkload, error) {
	assistantHours, err := uc.repository.SumWeeklyHoursByUserAndRole(assignment.UserID, domain.RoleAssistant)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	monitorHours, err := uc.repository.SumWeeklyHoursByUserAndRole(assignment.UserID, domain.RoleMonitor)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	monitorCount, err := uc.repository.CountAssignmentsByUserAndRole(assignment.UserID, domain.RoleMonitor)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	workload := domain.UserAssignmentWorkload{
		AssistantWeeklyHours: assistantHours,
		MonitorWeeklyHours:   monitorHours,
		MonitorAssignments:   monitorCount,
	}

	if assignment.Role == domain.RoleAssistant {
		workload.AssistantWeeklyHours -= assignment.WeeklyHours
	}

	if assignment.Role == domain.RoleMonitor {
		workload.MonitorWeeklyHours -= assignment.WeeklyHours
		workload.MonitorAssignments--
	}

	return workload, nil
}

func (uc *UpdateAssignment) Execute(input UpdateAssignmentInput) (*UpdateAssignmentOutput, error) {
	assignment, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.validateWorkspaceNotClosed(assignment.WorkspaceID); err != nil {
		return nil, err
	}

	if err := uc.ensureNoExactDuplicate(assignment, input.Role); err != nil {
		return nil, err
	}

	currentWorkload, err := uc.buildCurrentWorkload(assignment)
	if err != nil {
		return nil, err
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
