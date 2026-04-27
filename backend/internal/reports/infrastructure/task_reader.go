package infrastructure

import (
	tasksDomain "backend/internal/tasks/domain"
	"backend/internal/shared/database"
)

type TaskReader struct{}

func NewTaskReader() *TaskReader {
	return &TaskReader{}
}

func (r *TaskReader) FindAll() ([]tasksDomain.Task, error) {
	var tasks []tasksDomain.Task
	err := database.DB.Order("id asc").Find(&tasks).Error
	return tasks, err
}
