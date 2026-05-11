package domain

type TaskRepository interface {
	Create(task *Task) error
	FindByID(id uint) (*Task, error)
	FindAll() ([]Task, error)
	FindAllByUserID(userID uint) ([]Task, error)
	Update(task *Task) error
	UpdateAttachments(id uint, attachments []TaskAttachment) error
	Delete(id uint) error
}
