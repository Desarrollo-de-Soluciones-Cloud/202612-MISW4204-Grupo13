package infrastructure

import (
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/shared/database"
)

type AssignmentReader struct{}

func NewAssignmentReader() *AssignmentReader {
	return &AssignmentReader{}
}

func (r *AssignmentReader) FindAllByWorkspaceID(workspaceID uint) ([]assignmentsDomain.Assignment, error) {
	var assignments []assignmentsDomain.Assignment
	err := database.DB.Where("workspace_id = ?", workspaceID).Order("id asc").Find(&assignments).Error
	return assignments, err
}
