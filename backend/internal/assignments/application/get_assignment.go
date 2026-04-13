package application

import "backend/internal/assignments/domain"

type GetAssignmentByIDInput struct {
	ID uint
}

type GetAssignmentByIDOutput struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type GetAssignmentByID struct {
	repository domain.AssignmentRepository
}

func NewGetAssignmentByID(repo domain.AssignmentRepository) *GetAssignmentByID {
	return &GetAssignmentByID{repository: repo}
}

func (uc *GetAssignmentByID) Execute(input GetAssignmentByIDInput) (*GetAssignmentByIDOutput, error) {
	assignment, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetAssignmentByIDOutput{
		ID:          assignment.ID,
		UserID:      assignment.UserID,
		WorkspaceID: assignment.WorkspaceID,
		Role:        assignment.Role,
		WeeklyHours: assignment.WeeklyHours,
	}, nil
}
