package domain

type AssignmentRepository interface {
	Create(assignment *Assignment) error
	FindByID(id uint) (*Assignment, error)
	FindAll() ([]Assignment, error)
	FindAllByUserID(userID uint) ([]Assignment, error)
	FindByWorkspaceUserID(workspaceUserID uint) ([]Assignment, error)
	FindByWorkspaceIDsAndRoles(workspaceIDs []uint, roles []AssignmentRole) ([]Assignment, error)
	SumWeeklyHoursByUserAndRole(userID uint, role AssignmentRole) (int, error)
	CountAssignmentsByUserAndRole(userID uint, role AssignmentRole) (int, error)
	Update(assignment *Assignment) error
}
