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

	if input.UserID > 0 && input.UserID != workspace.UserID {
		return nil, workspacesDomain.ErrWorkspaceUserIDChangeNotAllowed
	}

	if err := uc.validateWorkspaceOwner(workspace.UserID); err != nil {
		return nil, err
	}

	periodID := workspace.PeriodID
	if input.PeriodID > 0 {
		periodID = input.PeriodID
	}

	period, err := uc.findOpenPeriod(periodID)
	if err != nil {
		return nil, err
	}
	workspace.PeriodID = periodID

	if err := validateWorkspaceDatesWithinPeriod(input, period); err != nil {
		return nil, err
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

func (uc *UpdateWorkspace) validateWorkspaceOwner(userID uint) error {
	user, err := uc.userRepository.FindByID(userID)
	if err != nil {
		return workspacesDomain.ErrWorkspaceUserNotFound
	}
	if user.GlobalRole != usersDomain.RoleProfessor {
		return workspacesDomain.ErrWorkspaceUserNotProfessor
	}
	return nil
}

func (uc *UpdateWorkspace) findOpenPeriod(periodID uint) (*periodsDomain.Period, error) {
	period, err := uc.periodRepository.FindByID(periodID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspacePeriodNotFound
	}
	if period.PeriodState == periodsDomain.ClosedPeriod {
		return nil, workspacesDomain.ErrWorkspacePeriodClosed
	}
	return period, nil
}

func validateWorkspaceDatesWithinPeriod(input UpdateWorkspaceInput, period *periodsDomain.Period) error {
	if input.InitialDate < period.InitialDate {
		return workspacesDomain.ErrWorkspaceInitialDateOutOfRange
	}
	if input.FinalDate > period.FinalDate {
		return workspacesDomain.ErrWorkspaceFinalDateOutOfRange
	}
	return nil
}
