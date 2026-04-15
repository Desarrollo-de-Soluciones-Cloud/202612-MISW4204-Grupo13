package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type TaskAssignmentRepository interface {
	FindByID(id uint) (*assignmentsDomain.Assignment, error)
}

type TaskWorkspaceRepository interface {
	FindByID(id uint) (*workspacesDomain.Workspace, error)
}

type TaskWeekRepository interface {
	FindByID(id uint) (*weeksDomain.Week, error)
}
