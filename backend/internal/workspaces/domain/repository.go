package domain

type WorkspaceRepository interface {
	Create(workspace *Workspace) error
	FindByID(id uint) (*Workspace, error)
	FindAll() ([]Workspace, error)
	FindByPeriodID(periodID uint) ([]Workspace, error)
	FindByUserID(userID uint) ([]Workspace, error)
	Update(workspace *Workspace) error
	Delete(id uint) error
}
