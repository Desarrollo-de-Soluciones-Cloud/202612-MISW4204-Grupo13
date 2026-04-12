package application

import "backend/internal/assignments/domain"

type UpdateAssignmentInput struct {
	ID          uint
	Role        domain.AssignmentRole
	WeeklyHours int
}

type UpdateAssignmentOutput struct {
	ID          uint                  `json:"id"`
	UserID      uint                  `json:"user_id"`
	WorkspaceID uint                  `json:"workspace_id"`
	Role        domain.AssignmentRole `json:"role"`
	WeeklyHours int                   `json:"weekly_hours"`
}

type UpdateAssignment struct {
	repository domain.AssignmentRepository
}

func NewUpdateAssignment(repo domain.AssignmentRepository) *UpdateAssignment {
	return &UpdateAssignment{repository: repo}
}

func (uc *UpdateAssignment) Execute(input UpdateAssignmentInput) (*UpdateAssignmentOutput, error) {
	assignment, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	if err := assignment.UpdateAdmin(input.Role, input.WeeklyHours); err != nil {
		return nil, err
	}

	if err := uc.repository.Update(assignment); err != nil {
		return nil, err
	}

	return &UpdateAssignmentOutput{
		ID:          assignment.ID,
		UserID:      assignment.UserID,
		WorkspaceID: assignment.WorkspaceID,
		Role:        assignment.Role,
		WeeklyHours: assignment.WeeklyHours,
	}, nil
}
