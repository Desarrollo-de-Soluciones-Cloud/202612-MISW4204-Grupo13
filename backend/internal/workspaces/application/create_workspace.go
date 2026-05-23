package application

import (
	periodsDomain "backend/internal/periods/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"time"
)




type CreateWorkspaceInput struct {
	PeriodID     uint
	UserID       uint
	Name         string
	Type         string
	InitialDate  string
	FinalDate    string
	Observations string
	State        string
}

type CreateWorkspaceOutput struct {
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

type CreateWorkspace struct {
	workspaceRepository workspacesDomain.WorkspaceRepository
	periodRepository    periodsDomain.PeriodRepository
	userRepository      usersDomain.UserRepository
}

func NewCreateWorkspace(workspaceRepo workspacesDomain.WorkspaceRepository, periodRepo periodsDomain.PeriodRepository, userRepo usersDomain.UserRepository) *CreateWorkspace {
	return &CreateWorkspace{
		workspaceRepository: workspaceRepo,
		periodRepository:    periodRepo,
		userRepository:      userRepo,
	}
}

func (uc *CreateWorkspace) Execute(input CreateWorkspaceInput) (*CreateWorkspaceOutput, error) {
	// Verify that workspace initial date is not greater than final date
	if input.InitialDate > input.FinalDate {
		return nil, workspacesDomain.ErrWorkspaceDateSequenceInvalid
	}

	// Verify that the period exists
	period, err := uc.periodRepository.FindByID(input.PeriodID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspacePeriodNotFound
	}
	if period.PeriodState == periodsDomain.ClosedPeriod {
		return nil, workspacesDomain.ErrWorkspacePeriodClosed
	}
	if period.InscriptionFinalDate < time.Now().Format("2006-01-02") {
		return nil, workspacesDomain.ErrWorkspaceInscriptionClosed
	}

	// Verify that workspace dates are within period date range
	if input.InitialDate < period.InitialDate {
		return nil, workspacesDomain.ErrWorkspaceInitialDateOutOfRange
	}
	if input.FinalDate > period.FinalDate {
		return nil, workspacesDomain.ErrWorkspaceFinalDateOutOfRange
	}

	// Verify that the user exists and has professor role
	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		return nil, workspacesDomain.ErrWorkspaceUserNotFound
	}
	if user.GlobalRole != usersDomain.RoleProfessor {
		return nil, workspacesDomain.ErrWorkspaceUserNotProfessor
	}

	workspace, err := workspacesDomain.NewWorkspace(workspacesDomain.WorkspaceInput{
		PeriodID:     input.PeriodID,
		UserID:       input.UserID,
		Name:         input.Name,
		Type:         workspacesDomain.WorkspaceType(input.Type),
		InitialDate:  input.InitialDate,
		FinalDate:    input.FinalDate,
		Observations: input.Observations,
		State:        workspacesDomain.WorkspaceState(input.State),
	})
	if err != nil {
		return nil, err
	}

	if err := uc.workspaceRepository.Create(workspace); err != nil {
		return nil, err
	}

	return &CreateWorkspaceOutput{
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
