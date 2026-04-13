package domain

import (
	"strings"
	"time"
)

func ValidateTaskAssignmentID(assignmentID uint) error {
	if assignmentID == 0 {
		return ErrTaskAssignmentIDRequired
	}
	return nil
}

func NormalizeTaskTitle(title string) string {
	return strings.TrimSpace(title)
}

func NormalizeTaskDescription(description string) string {
	return strings.TrimSpace(description)
}

func NormalizeTaskObservations(observations string) string {
	return strings.TrimSpace(observations)
}

func ValidateTaskTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrTaskTitleRequired
	}
	return nil
}

func ValidateTaskDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return ErrTaskDescriptionRequired
	}
	return nil
}

func ValidateTaskStatus(status TaskStatus) error {
	switch {
	case strings.TrimSpace(string(status)) == "":
		return ErrTaskStatusRequired
	case !IsValidTaskStatus(status):
		return ErrTaskStatusInvalid
	default:
		return nil
	}
}

func ValidateTaskSpentHours(spentHours int) error {
	switch {
	case spentHours == 0:
		return ErrTaskSpentHoursRequired
	case spentHours < 1:
		return ErrTaskSpentHoursInvalid
	default:
		return nil
	}
}

func ValidateTaskWeekStartDate(weekStartDate time.Time) error {
	if weekStartDate.IsZero() {
		return ErrTaskWeekStartDateRequired
	}
	return nil
}

func NormalizeWeekStartDate(date time.Time) time.Time {
	if date.IsZero() {
		return time.Time{}
	}

	normalizedDate := normalizeDateOnly(date)
	weekday := int(normalizedDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return normalizedDate.AddDate(0, 0, -(weekday - 1))
}

func IsWeekClosed(weekStartDate, now time.Time) bool {
	normalizedWeekStart := NormalizeWeekStartDate(weekStartDate)
	if normalizedWeekStart.IsZero() {
		return false
	}

	return NormalizeWeekStartDate(now).After(normalizedWeekStart)
}

func IsWeekActive(weekStartDate, now time.Time) bool {
	normalizedWeekStart := NormalizeWeekStartDate(weekStartDate)
	if normalizedWeekStart.IsZero() {
		return false
	}

	return NormalizeWeekStartDate(now).Equal(normalizedWeekStart)
}

func normalizeDateOnly(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
