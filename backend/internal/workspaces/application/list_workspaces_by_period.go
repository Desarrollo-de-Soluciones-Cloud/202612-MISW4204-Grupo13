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

	return &ListWorkspacesByPeriodOutput{Workspaces: mapWorkspaceDTOs(workspaces)}, nil
}
