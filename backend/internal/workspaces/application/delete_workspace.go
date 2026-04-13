package application

import (
	"backend/internal/workspaces/domain"
)

type DeleteWorkspaceInput struct {
	ID uint
}

type DeleteWorkspace struct {
	repository domain.WorkspaceRepository
}

func NewDeleteWorkspace(repo domain.WorkspaceRepository) *DeleteWorkspace {
	return &DeleteWorkspace{repository: repo}
}

func (uc *DeleteWorkspace) Execute(input DeleteWorkspaceInput) error {
	return uc.repository.Delete(input.ID)
}
