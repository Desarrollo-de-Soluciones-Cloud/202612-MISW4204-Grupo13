package domain

type AssignmentRole string

const (
	RoleMonitor   AssignmentRole = "monitor"
	RoleAssistant AssignmentRole = "assistant"

	MaxAssistantWeeklyHours = 22
	MaxMonitorAssignments   = 3
	MaxMonitorWeeklyHours   = 12
	MonitorHoursPercentage  = 40
)

func IsValidAssignmentRole(role AssignmentRole) bool {
	switch role {
	case RoleMonitor, RoleAssistant:
		return true
	default:
		return false
	}
}
