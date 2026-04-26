package application

import (
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type UpdateWorkspaceInput struct {
	ID           uint
	PeriodID     uint
	UserID       uint
	Name         string
	Type         string
	InitialDate  string
	FinalDate    string
	Observations string
	State        string
}

type UpdateWorkspaceOutput struct {
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

type UpdateWorkspace struct {
	workspaceRepository workspacesDomain.WorkspaceRepository
	periodRepository    periodsDomain.PeriodRepository
	userRepository      usersDomain.UserRepository
}

func NewUpdateWorkspace(workspaceRepo workspacesDomain.WorkspaceRepository, periodRepo periodsDomain.PeriodRepository, userRepo usersDomain.UserRepository) *UpdateWorkspace {
	return &UpdateWorkspace{
		workspaceRepository: workspaceRepo,
		periodRepository:    periodRepo,
		userRepository:      userRepo,
	}
}

func (uc *UpdateWorkspace) Execute(input UpdateWorkspaceInput) (*UpdateWorkspaceOutput, error) {
	// Verify that workspace initial date is not greater than final date
	if input.InitialDate > input.FinalDate {
		return nil, workspacesDomain.ErrWorkspaceDateSequenceInvalid
	}

	workspace, err := uc.workspaceRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if workspace.State == workspacesDomain.ClosedState {
		return nil, workspacesDomain.ErrWorkspaceClosedUpdateForbidden
	}

	// Verify that user_id cannot be changed
	if input.UserID > 0 && input.UserID != workspace.UserID {
		return nil, workspacesDomain.ErrWorkspaceUserIDChangeNotAllowed
	}

	// Verify that the user still has professor role (validate existing user)
	user, err := uc.userRepository.FindByID(workspace.UserID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspaceUserNotFound
	}
	if user.GlobalRole != usersDomain.RoleProfessor {
		return nil, workspacesDomain.ErrWorkspaceUserNotProfessor
	}

	// Verify that the period exists if PeriodID is provided
	if input.PeriodID > 0 {
		period, err := uc.periodRepository.FindByID(input.PeriodID)
		if err != nil {
			return nil, workspacesDomain.ErrWorkspacePeriodNotFound
		}
		if period.PeriodState == periodsDomain.ClosedPeriod {
			return nil, workspacesDomain.ErrWorkspacePeriodClosed
		}
		workspace.PeriodID = input.PeriodID
	}

	// Get the period to validate workspace dates are within period range
	period, err := uc.periodRepository.FindByID(workspace.PeriodID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspacePeriodNotFound
	}
	if period.PeriodState == periodsDomain.ClosedPeriod {
		return nil, workspacesDomain.ErrWorkspacePeriodClosed
	}

	// Verify that workspace dates are within period date range
	if input.InitialDate < period.InitialDate {
		return nil, workspacesDomain.ErrWorkspaceInitialDateOutOfRange
	}
	if input.FinalDate > period.FinalDate {
		return nil, workspacesDomain.ErrWorkspaceFinalDateOutOfRange
	}

	if err := workspace.UpdateWorkspace(
		input.Name,
		workspacesDomain.WorkspaceType(input.Type),
		input.InitialDate,
		input.FinalDate,
		input.Observations,
		workspacesDomain.WorkspaceState(input.State),
	); err != nil {
		return nil, err
	}

	if err := uc.workspaceRepository.Update(workspace); err != nil {
		return nil, err
	}

	return &UpdateWorkspaceOutput{
		ID:           workspace.ID,
		PeriodID:     workspace.PeriodID,
		UserID:       workspace.UserID,
		Name:         workspace.Name,
		Type:         string(workspace.Type),
		InitialDate:  workspace.InitialDate,
		FinalDate:    workspace.FinalDate,
		Observations: workspace.Observations,
		State:        string(workspace.State),
	}, nil
}
