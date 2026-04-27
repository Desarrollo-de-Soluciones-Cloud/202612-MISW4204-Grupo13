package application

import (
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type CloseWorkspaceInput struct {
	ID uint
}

type CloseWorkspaceOutput struct {
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

type CloseWorkspace struct {
	workspaceRepository workspacesDomain.WorkspaceRepository
	userRepository      usersDomain.UserRepository
}

func NewCloseWorkspace(workspaceRepo workspacesDomain.WorkspaceRepository, userRepo usersDomain.UserRepository) *CloseWorkspace {
	return &CloseWorkspace{
		workspaceRepository: workspaceRepo,
		userRepository:      userRepo,
	}
}

func (uc *CloseWorkspace) Execute(input CloseWorkspaceInput) (*CloseWorkspaceOutput, error) {
	workspace, err := uc.workspaceRepository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	// Verify that the user still has professor role
	user, err := uc.userRepository.FindByID(workspace.UserID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspaceUserNotFound
	}
	if user.GlobalRole != usersDomain.RoleProfessor {
		return nil, workspacesDomain.ErrWorkspaceUserNotProfessor
	}

	// Update workspace state to closed
	workspace.State = workspacesDomain.ClosedState

	if err := uc.workspaceRepository.Update(workspace); err != nil {
		return nil, err
	}

	return &CloseWorkspaceOutput{
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
