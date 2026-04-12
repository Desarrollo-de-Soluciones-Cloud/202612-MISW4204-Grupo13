package domain

type AssignmentRepository interface {
	Create(assignment *Assignment) error
	FindByID(id uint) (*Assignment, error)
	FindAllByUserID(userID uint) ([]Assignment, error)
	Update(assignment *Assignment) error
}
