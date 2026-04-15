package infrastructure

import (
	"errors"

	"backend/internal/shared/database"
	"backend/internal/tasks/domain"

	"gorm.io/gorm"
)

type TaskRepository struct{}

var _ domain.TaskRepository = (*TaskRepository)(nil)

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Attachments").Create(task).Error; err != nil {
			return err
		}
		if len(task.Attachments) == 0 {
			return nil
		}
		for i := range task.Attachments {
			task.Attachments[i].TaskID = task.ID
		}
		return tx.Create(&task.Attachments).Error
	})
}

func (r *TaskRepository) FindByID(id uint) (*domain.Task, error) {
	var task domain.Task
	result := database.DB.Preload("Attachments").First(&task, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, result.Error
	}
	return &task, nil
}

func (r *TaskRepository) FindAll() ([]domain.Task, error) {
	var tasks []domain.Task
	result := database.DB.Preload("Attachments").Order("id asc").Find(&tasks)
	return tasks, result.Error
}

func (r *TaskRepository) FindAllByUserID(userID uint) ([]domain.Task, error) {
	var tasks []domain.Task
	result := database.DB.Preload("Attachments").Where("user_id = ?", userID).Order("id asc").Find(&tasks)
	return tasks, result.Error
}

func (r *TaskRepository) Update(task *domain.Task) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Attachments").Save(task).Error; err != nil {
			return err
		}
		if err := tx.Where("task_id = ?", task.ID).Delete(&domain.TaskAttachment{}).Error; err != nil {
			return err
		}
		if len(task.Attachments) == 0 {
			return nil
		}
		for i := range task.Attachments {
			task.Attachments[i].ID = 0
			task.Attachments[i].TaskID = task.ID
		}
		return tx.Create(&task.Attachments).Error
	})
}

func (r *TaskRepository) Delete(id uint) error {
	result := database.DB.Delete(&domain.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Task{}, &domain.TaskAttachment{})
}

func (r *TaskRepository) NormalizeLegacyStatuses() error {
	legacyStatuses := map[string]domain.TaskStatus{
		"open":           domain.TaskStatusAbierto,
		"in_development": domain.TaskStatusEnDesarrollo,
		"finished":       domain.TaskStatusFinalizado,
	}

	for legacy, normalized := range legacyStatuses {
		if err := database.DB.Model(&domain.Task{}).
			Where("status = ?", legacy).
			Update("status", string(normalized)).Error; err != nil {
			return err
		}
	}

	return nil
}
