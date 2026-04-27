package infrastructure

import (
	"errors"

	"backend/internal/shared/database"
	"backend/internal/workspaces/domain"

	"gorm.io/gorm"
)

type WorkspaceRepository struct{}

var _ domain.WorkspaceRepository = (*WorkspaceRepository)(nil)

func NewWorkspaceRepository() *WorkspaceRepository {
	return &WorkspaceRepository{}
}

func (r *WorkspaceRepository) Create(workspace *domain.Workspace) error {
	result := database.DB.Create(workspace)
	return result.Error
}

func (r *WorkspaceRepository) FindByID(id uint) (*domain.Workspace, error) {
	var workspace domain.Workspace
	result := database.DB.First(&workspace, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkspaceNotFound
		}
		return nil, result.Error
	}
	return &workspace, nil
}

func (r *WorkspaceRepository) FindAll() ([]domain.Workspace, error) {
	var workspaces []domain.Workspace
	result := database.DB.Find(&workspaces)
	return workspaces, result.Error
}

func (r *WorkspaceRepository) FindByPeriodID(periodID uint) ([]domain.Workspace, error) {
	var workspaces []domain.Workspace
	result := database.DB.Where("period_id = ?", periodID).Find(&workspaces)
	return workspaces, result.Error
}

func (r *WorkspaceRepository) FindByUserID(userID uint) ([]domain.Workspace, error) {
	var workspaces []domain.Workspace
	result := database.DB.Where("user_id = ?", userID).Find(&workspaces)
	return workspaces, result.Error
}

func (r *WorkspaceRepository) Update(workspace *domain.Workspace) error {
	return database.DB.Save(workspace).Error
}

func (r *WorkspaceRepository) Delete(id uint) error {
	result := database.DB.Delete(&domain.Workspace{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrWorkspaceNotFound
	}
	return nil
}

func (r *WorkspaceRepository) AutoMigrate() error {
	return database.DB.AutoMigrate(&domain.Workspace{})
}
