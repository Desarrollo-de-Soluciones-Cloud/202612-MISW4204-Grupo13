package application

import (
	"backend/internal/workspaces/domain"
)

type GetWorkspaceByIDInput struct {
	ID uint
}

type GetWorkspaceByIDOutput struct {
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

type GetWorkspaceByID struct {
	repository domain.WorkspaceRepository
}

func NewGetWorkspaceByID(repo domain.WorkspaceRepository) *GetWorkspaceByID {
	return &GetWorkspaceByID{repository: repo}
}

func (uc *GetWorkspaceByID) Execute(input GetWorkspaceByIDInput) (*GetWorkspaceByIDOutput, error) {
	workspace, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetWorkspaceByIDOutput{
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
