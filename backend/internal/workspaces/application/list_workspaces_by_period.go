package application

import (
	"backend/internal/workspaces/domain"
)

type ListWorkspacesByPeriodInput struct {
	PeriodID uint
}

type ListWorkspacesByPeriodOutput struct {
	Workspaces []WorkspaceDTO `json:"workspaces"`
}

type ListWorkspacesByPeriod struct {
	repository domain.WorkspaceRepository
}

func NewListWorkspacesByPeriod(repo domain.WorkspaceRepository) *ListWorkspacesByPeriod {
	return &ListWorkspacesByPeriod{repository: repo}
}

func (uc *ListWorkspacesByPeriod) Execute(input ListWorkspacesByPeriodInput) (*ListWorkspacesByPeriodOutput, error) {
	workspaces, err := uc.repository.FindByPeriodID(input.PeriodID)
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

	return &ListWorkspacesByPeriodOutput{Workspaces: result}, nil
}
