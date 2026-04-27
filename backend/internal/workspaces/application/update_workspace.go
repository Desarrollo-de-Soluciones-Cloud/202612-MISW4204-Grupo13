package application

import (
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type UpdateWorkspaceInput struct {
	ID           uint
	PeriodID     uint
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
	workspace, err := uc.workspaceRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}
	if workspace.State == workspacesDomain.ClosedState {
		return nil, workspacesDomain.ErrWorkspaceClosedUpdateForbidden
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
		_, err := uc.periodRepository.FindByID(input.PeriodID)
		if err != nil {
			return nil, workspacesDomain.ErrWorkspacePeriodNotFound
		}
		workspace.PeriodID = input.PeriodID
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
