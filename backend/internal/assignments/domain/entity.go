package domain

import "time"

type Assignment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	WorkspaceID uint           `gorm:"not null;index" json:"workspace_id"`
	Role        AssignmentRole `gorm:"size:20;not null" json:"role"`
	WeeklyHours int            `gorm:"not null" json:"weekly_hours"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (a *Assignment) UpdateAdmin(role AssignmentRole, weeklyHours int) error {
	if err := ValidateAssignmentRole(role); err != nil {
		return err
	}
	if err := ValidateAssignmentWeeklyHours(weeklyHours); err != nil {
		return err
	}

	a.Role = role
	a.WeeklyHours = weeklyHours

	return nil
}

func NewAssignment(userID, workspaceID uint, role AssignmentRole, weeklyHours int) (*Assignment, error) {
	if err := ValidateAssignmentUserID(userID); err != nil {
		return nil, err
	}
	if err := ValidateAssignmentWorkspaceID(workspaceID); err != nil {
		return nil, err
	}
	if err := ValidateAssignmentRole(role); err != nil {
		return nil, err
	}
	if err := ValidateAssignmentWeeklyHours(weeklyHours); err != nil {
		return nil, err
	}

	return &Assignment{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
		WeeklyHours: weeklyHours,
	}, nil
}
