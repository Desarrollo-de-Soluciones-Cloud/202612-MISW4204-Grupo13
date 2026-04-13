package domain

type UserAssignmentWorkload struct {
	AssistantWeeklyHours int
	MonitorWeeklyHours   int
	MonitorAssignments   int
}

func BuildWorkloadWithAssignment(current UserAssignmentWorkload, role AssignmentRole, weeklyHours int) (UserAssignmentWorkload, error) {
	if err := ValidateAssignmentRole(role); err != nil {
		return UserAssignmentWorkload{}, err
	}
	if err := ValidateAssignmentWeeklyHours(weeklyHours); err != nil {
		return UserAssignmentWorkload{}, err
	}

	next := current
	switch role {
	case RoleAssistant:
		next.AssistantWeeklyHours += weeklyHours
	case RoleMonitor:
		next.MonitorWeeklyHours += weeklyHours
		next.MonitorAssignments++
	}

	return next, nil
}

func ValidateRF05Workload(workload UserAssignmentWorkload) error {
	if workload.AssistantWeeklyHours > MaxAssistantWeeklyHours {
		return ErrAssignmentAssistantHoursLimitExceeded
	}

	if workload.MonitorAssignments > MaxMonitorAssignments {
		return ErrAssignmentMonitorCountLimitExceeded
	}

	if workload.MonitorWeeklyHours > MaxMonitorWeeklyHours {
		return ErrAssignmentMonitorHoursLimitExceeded
	}

	if workload.AssistantWeeklyHours > 0 {
		maxMonitorHours := CalculateMaxMonitorHoursFromAssistant(workload.AssistantWeeklyHours)
		if workload.MonitorWeeklyHours > maxMonitorHours {
			return ErrAssignmentMonitorFortyPercentExceeded
		}
	}

	return nil
}

func CalculateMaxMonitorHoursFromAssistant(assistantWeeklyHours int) int {
	if assistantWeeklyHours <= 0 {
		return 0
	}

	return (assistantWeeklyHours*MonitorHoursPercentage + 99) / 100
}
