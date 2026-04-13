package domain

type TaskStatus string

const (
	TaskStatusOpen          TaskStatus = "open"
	TaskStatusInDevelopment TaskStatus = "in_development"
	TaskStatusFinished      TaskStatus = "finished"
)

func IsValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusOpen, TaskStatusInDevelopment, TaskStatusFinished:
		return true
	default:
		return false
	}
}
