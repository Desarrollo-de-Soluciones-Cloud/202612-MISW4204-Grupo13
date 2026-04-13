package domain

type TaskRepository interface {
	Create(task *Task) error
	FindByID(id uint) (*Task, error)
	FindAll() ([]Task, error)
	FindAllByUserID(userID uint) ([]Task, error)
	Update(task *Task) error
	Delete(id uint) error
}
