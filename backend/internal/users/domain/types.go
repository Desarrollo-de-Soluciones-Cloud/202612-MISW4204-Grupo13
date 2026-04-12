package domain

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleProfessor UserRole = "professor"
	RoleMonitor   UserRole = "monitor"
	RoleAssistant UserRole = "assistant"
)

func IsValidUserRole(role UserRole) bool {
	switch role {
	case RoleAdmin, RoleProfessor, RoleMonitor, RoleAssistant:
		return true
	default:
		return false
	}
}
