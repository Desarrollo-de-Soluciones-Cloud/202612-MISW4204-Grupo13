package infrastructure

import (
	"backend/internal/shared/database"
	tasksDomain "backend/internal/tasks/domain"
)

type TaskReader struct{}

func NewTaskReader() *TaskReader {
	return &TaskReader{}
}

func (r *TaskReader) FindAllByWorkspaceAndWeek(
	workspaceID uint,
	weekID uint,
	weekInitialDate string,
) ([]tasksDomain.Task, error) {
	var tasks []tasksDomain.Task

	err := database.DB.
		Joins("JOIN assignments ON assignments.id = tasks.assignment_id").
		Where("assignments.workspace_id = ?", workspaceID).
		Where("(tasks.week_id = ? OR tasks.week_start_date = ?)", weekID, weekInitialDate).
		Order("tasks.id asc").
		Find(&tasks).
		Error

	return tasks, err
}