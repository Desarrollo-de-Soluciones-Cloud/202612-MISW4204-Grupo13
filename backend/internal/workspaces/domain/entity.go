package domain

import (
	"time"
)

type Workspace struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	PeriodID	 uint           `gorm:"not null" json:"period_id"`
	UserID	     uint           `gorm:"not null" json:"user_id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Type         WorkspaceType  `gorm:"size:100;not null" json:"type"`
	InitialDate  string         `gorm:"size:100;not null" json:"initial_date"`
	FinalDate    string         `gorm:"size:100;not null" json:"final_date"`
	Observations string         `gorm:"type:text;not null" json:"inscription_final_date"`
	State        WorkspaceState `gorm:"size:20;not null" json:"period_state"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func NewWorkspace(periodID, userID uint, name string, workspaceType WorkspaceType, initialDate, finalDate, observations string, state WorkspaceState) (*Workspace, error) {
	normalizedName := NormalizeWorkspaceName(name)

	if err := ValidateWorkspaceName(normalizedName); err != nil {
		return nil, err
	}
	if err := ValidateWorkspacePeriodID(periodID); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceUserID(userID); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceType(workspaceType); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceInitialDate(initialDate); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceFinalDate(finalDate); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceDateSequence(initialDate, finalDate); err != nil {
		return nil, err
	}
	if err := ValidateWorkspaceState(state); err != nil {
		return nil, err
	}

	return &Workspace{
		PeriodID:     periodID,
		UserID:       userID,
		Name:         normalizedName,
		Type:         workspaceType,
		InitialDate:  initialDate,
		FinalDate:    finalDate,
		Observations: observations,
		State:        state,
	}, nil
}

func (w *Workspace) UpdateWorkspace(name string, workspaceType WorkspaceType, initialDate, finalDate, observations string, state WorkspaceState) error {
	normalizedName := NormalizeWorkspaceName(name)

	if err := ValidateWorkspaceName(normalizedName); err != nil {
		return err
	}
	if err := ValidateWorkspaceType(workspaceType); err != nil {
		return err
	}
	if err := ValidateWorkspaceInitialDate(initialDate); err != nil {
		return err
	}
	if err := ValidateWorkspaceFinalDate(finalDate); err != nil {
		return err
	}
	if err := ValidateWorkspaceDateSequence(initialDate, finalDate); err != nil {
		return err
	}
	if err := ValidateWorkspaceState(state); err != nil {
		return err
	}

	w.Name = normalizedName
	w.Type = workspaceType
	w.InitialDate = initialDate
	w.FinalDate = finalDate
	w.Observations = observations
	w.State = state

	return nil
}