package domain

type TaskStatus string

const (
	TaskStatusAbierto      TaskStatus = "abierto"
	TaskStatusEnDesarrollo TaskStatus = "en desarrollo"
	TaskStatusFinalizado   TaskStatus = "finalizado"
)

func IsValidTaskStatus(status TaskStatus) bool {
	switch status {
	case TaskStatusAbierto, TaskStatusEnDesarrollo, TaskStatusFinalizado:
		return true
	default:
		return false
	}
}
