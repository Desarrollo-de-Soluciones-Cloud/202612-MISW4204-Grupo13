package application

import (
	"backend/internal/assignments/domain"
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
)

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
	repository          domain.AssignmentRepository
	userRepository      AssignmentUserRepository
	workspaceRepository AssignmentWorkspaceRepository
	periodRepository    AssignmentPeriodRepository
}

type AssignmentUserRepository interface {
	FindByID(id uint) (*usersDomain.User, error)
}

type AssignmentWorkspaceRepository interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type AssignmentPeriodRepository interface {
	FindByID(id uint) (*periodsDomain.Period, error)
}

func NewCreateAssignment(repo domain.AssignmentRepository) *CreateAssignment {
	return &CreateAssignment{repository: repo}
}

func (uc *CreateAssignment) WithRepositories(userRepo AssignmentUserRepository, workspaceRepo AssignmentWorkspaceRepository) *CreateAssignment {
	uc.userRepository = userRepo
	uc.workspaceRepository = workspaceRepo

	return uc
}

func (uc *CreateAssignment) WithPeriodRepository(periodRepo AssignmentPeriodRepository) *CreateAssignment {
	uc.periodRepository = periodRepo
	return uc
}

func mapWorkspaceNotFoundError(err error) bool {
	return errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) || err.Error() == workspacesDomain.ErrWorkspaceNotFound.Error()
}

func (uc *CreateAssignment) validateUserAndWorkspace(userID, workspaceID uint) error {
	if uc.userRepository != nil {
		if _, err := uc.userRepository.FindByID(userID); err != nil {
			if errors.Is(err, usersDomain.ErrUserNotFound) {
				return domain.ErrAssignmentUserNotFound
			}
			return err
		}
	}

	if uc.workspaceRepository != nil {
		workspace, err := uc.workspaceRepository.FindByID(workspaceID)
		if err != nil {
			if mapWorkspaceNotFoundError(err) {
				return domain.ErrAssignmentWorkspaceNotFound
			}
			return err
		}
		if workspace.State == workspacesDomain.ClosedState {
			return domain.ErrAssignmentWorkspaceClosed
		}

		// Validate period is not closed
		if uc.periodRepository != nil {
			period, err := uc.periodRepository.FindByID(workspace.PeriodID)
			if err != nil {
				if errors.Is(err, periodsDomain.ErrPeriodNotFound) {
					return domain.ErrAssignmentWorkspaceNotFound
				}
				return err
			}
			if period.PeriodState == periodsDomain.ClosedPeriod {
				return domain.ErrAssignmentPeriodClosed
			}
		}
	}

	return nil
}

func (uc *CreateAssignment) ensureNoExactDuplicate(userID, workspaceID uint, role domain.AssignmentRole) error {
	assignments, err := uc.repository.FindAllByUserID(userID)
	if err != nil {
		return err
	}

	for _, existing := range assignments {
		if existing.WorkspaceID == workspaceID && existing.Role == role {
			return domain.ErrAssignmentAlreadyExists
		}
	}

	return nil
}

func (uc *CreateAssignment) buildNextWorkload(userID uint, role domain.AssignmentRole, weeklyHours int) (domain.UserAssignmentWorkload, error) {
	assistantHours, err := uc.repository.SumWeeklyHoursByUserAndRole(userID, domain.RoleAssistant)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	monitorHours, err := uc.repository.SumWeeklyHoursByUserAndRole(userID, domain.RoleMonitor)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	monitorCount, err := uc.repository.CountAssignmentsByUserAndRole(userID, domain.RoleMonitor)
	if err != nil {
		return domain.UserAssignmentWorkload{}, err
	}

	return domain.BuildWorkloadWithAssignment(domain.UserAssignmentWorkload{
		AssistantWeeklyHours: assistantHours,
		MonitorWeeklyHours:   monitorHours,
		MonitorAssignments:   monitorCount,
	}, role, weeklyHours)
}

func (uc *CreateAssignment) Execute(input CreateAssignmentInput) (*CreateAssignmentOutput, error) {
	//nolint:godox // TODO RF04: Validar que user_id exista realmente cuando se defina la integracion con users.
	//nolint:godox // TODO RF04: Validar que workspace_id exista realmente cuando el modulo de workspaces este terminado.
	// TODO RF05: Validar RF05 usando solo vinculaciones activas cuando workspaces exponga estados activos/cerrados.
	// TODO RF05: Confirmar si RF05 debe ejecutarse con filtros por periodo academico cuando periodos y workspaces esten integrados.
	if err := uc.validateUserAndWorkspace(input.UserID, input.WorkspaceID); err != nil {
		return nil, err
	}

	if err := uc.ensureNoExactDuplicate(input.UserID, input.WorkspaceID, input.Role); err != nil {
		return nil, err
	}

	nextWorkload, err := uc.buildNextWorkload(input.UserID, input.Role, input.WeeklyHours)
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
