package application

import (
	"backend/internal/workspaces/domain"
)

type ListWorkspacesOutput struct {
	Workspaces []WorkspaceDTO `json:"workspaces"`
}

type WorkspaceDTO struct {
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

type ListWorkspaces struct {
	repository domain.WorkspaceRepository
}

func NewListWorkspaces(repo domain.WorkspaceRepository) *ListWorkspaces {
	return &ListWorkspaces{repository: repo}
}

func (uc *ListWorkspaces) Execute() (*ListWorkspacesOutput, error) {
	workspaces, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	result := make([]WorkspaceDTO, len(workspaces))
	for i, w := range workspaces {
		result[i] = WorkspaceDTO{
			ID:           w.ID,
			PeriodID:     w.PeriodID,
			UserID:       w.UserID,
			Name:         w.Name,
			Type:         string(w.Type),
			InitialDate:  w.InitialDate,
			FinalDate:    w.FinalDate,
			Observations: w.Observations,
			State:        string(w.State),
		}
	}

	return &ListWorkspacesOutput{Workspaces: result}, nil
}
