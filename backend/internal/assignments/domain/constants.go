package domain

type AssignmentRole string

const (
	RoleMonitor   AssignmentRole = "monitor"
	RoleAssistant AssignmentRole = "assistant"
)

func IsValidAssignmentRole(role AssignmentRole) bool {
	switch role {
	case RoleMonitor, RoleAssistant:
		return true
	default:
		return false
	}
}
