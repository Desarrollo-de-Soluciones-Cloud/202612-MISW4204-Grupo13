package application

import "backend/internal/assignments/domain"

type ListAssignmentsByWorkspaceInput struct {
	ProfessorID uint
}

type ListAssignmentsByWorkspaceOutput struct {
	Assignments []AssignmentDTO `json:"assignments"`
}

type ListAssignmentsByWorkspace struct {
	repository domain.AssignmentRepository
}

func NewListAssignmentsByWorkspace(repo domain.AssignmentRepository) *ListAssignmentsByWorkspace {
	return &ListAssignmentsByWorkspace{repository: repo}
}

func (uc *ListAssignmentsByWorkspace) Execute(input ListAssignmentsByWorkspaceInput) (*ListAssignmentsByWorkspaceOutput, error) {
	if err := domain.ValidateAssignmentUserID(input.ProfessorID); err != nil {
		return nil, err
	}

	assignments, err := uc.repository.FindByWorkspaceUserID(input.ProfessorID)
	if err != nil {
		return nil, err
	}

	result := make([]AssignmentDTO, len(assignments))
	for i, a := range assignments {
		result[i] = AssignmentDTO{
			ID:          a.ID,
			UserID:      a.UserID,
			WorkspaceID: a.WorkspaceID,
			Role:        a.Role,
			WeeklyHours: a.WeeklyHours,
		}
	}

	return &ListAssignmentsByWorkspaceOutput{Assignments: result}, nil
}
