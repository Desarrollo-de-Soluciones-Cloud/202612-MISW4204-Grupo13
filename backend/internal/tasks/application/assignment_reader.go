package application

import assignmentsDomain "backend/internal/assignments/domain"

type TaskAssignmentRepository interface {
	FindByID(id uint) (*assignmentsDomain.Assignment, error)
}
