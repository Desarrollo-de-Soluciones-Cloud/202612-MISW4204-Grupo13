package application

import "backend/internal/assignments/domain"

type ListAssignmentsByUserIDInput struct {
	UserID uint
}

type AssignmentDTO struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type ListAssignmentsByUserIDOutput struct {
	Assignments []AssignmentDTO `json:"assignments"`
}

type ListAssignmentsByUserID struct {
	repository domain.AssignmentRepository
}

func NewListAssignmentsByUserID(repo domain.AssignmentRepository) *ListAssignmentsByUserID {
	return &ListAssignmentsByUserID{repository: repo}
}

func (uc *ListAssignmentsByUserID) Execute(input ListAssignmentsByUserIDInput) (*ListAssignmentsByUserIDOutput, error) {
	if err := domain.ValidateAssignmentUserID(input.UserID); err != nil {
		return nil, err
	}
	assignments, err := uc.repository.FindAllByUserID(input.UserID)
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

	return &ListAssignmentsByUserIDOutput{Assignments: result}, nil
}
