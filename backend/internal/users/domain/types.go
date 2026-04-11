package domain

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleProfessor UserRole = "professor"
	RoleStaff     UserRole = "staff"
)

func IsValidUserRole(role UserRole) bool {
	switch role {
	case RoleAdmin, RoleProfessor, RoleStaff:
		return true
	default:
		return false
	}
}
