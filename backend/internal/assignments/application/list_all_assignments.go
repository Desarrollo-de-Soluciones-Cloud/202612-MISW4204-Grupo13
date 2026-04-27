package application

import "backend/internal/assignments/domain"

type ListAllAssignmentsInput struct{}

type ListAllAssignmentsOutput struct {
	Assignments []AssignmentDTO `json:"assignments"`
}

type ListAllAssignments struct {
	repository domain.AssignmentRepository
}

func NewListAllAssignments(repo domain.AssignmentRepository) *ListAllAssignments {
	return &ListAllAssignments{repository: repo}
}

func (uc *ListAllAssignments) Execute(input ListAllAssignmentsInput) (*ListAllAssignmentsOutput, error) {
	assignments, err := uc.repository.FindAll()
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

	return &ListAllAssignmentsOutput{Assignments: result}, nil
}
