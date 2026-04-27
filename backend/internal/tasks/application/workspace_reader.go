package application

import workspacesDomain "backend/internal/workspaces/domain"

type TaskWorkspaceRepository interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}