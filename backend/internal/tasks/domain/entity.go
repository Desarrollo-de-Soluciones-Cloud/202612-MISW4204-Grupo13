package domain

import "time"

type Task struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null;index" json:"user_id"`
	AssignmentID  uint       `gorm:"not null;index" json:"assignment_id"`
	WeekID        *uint      `gorm:"default:null" json:"week_id"`
	Title         string     `gorm:"size:255;not null" json:"title"`
	Description   string     `gorm:"type:text;not null" json:"description"`
	Status        TaskStatus `gorm:"size:30;not null" json:"status"`
	SpentHours    int        `gorm:"not null" json:"spent_hours"`
	Observations  string     `gorm:"type:text;not null" json:"observations"`
	WeekStartDate time.Time  `gorm:"type:date;not null" json:"week_start_date"`
	Late          bool       `gorm:"not null;default:false" json:"late"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func NewTask(
	userID uint,
	assignmentID uint,
	weekID *uint,
	title, description string,
	status TaskStatus,
	spentHours int,
	observations string,
	weekStartDate time.Time,
	late bool,
) (*Task, error) {
	normalizedTitle := NormalizeTaskTitle(title)
	normalizedDescription := NormalizeTaskDescription(description)
	normalizedStatus := NormalizeTaskStatus(status)
	normalizedObservations := NormalizeTaskObservations(observations)
	normalizedWeekStartDate := NormalizeWeekStartDate(weekStartDate)

	if err := ValidateTaskAssignmentID(assignmentID); err != nil {
		return nil, err
	}
	if err := ValidateTaskTitle(normalizedTitle); err != nil {
		return nil, err
	}
	if err := ValidateTaskDescription(normalizedDescription); err != nil {
		return nil, err
	}
	if err := ValidateTaskStatus(normalizedStatus); err != nil {
		return nil, err
	}
	if err := ValidateTaskSpentHours(spentHours); err != nil {
		return nil, err
	}
	if err := ValidateTaskWeekStartDate(normalizedWeekStartDate); err != nil {
		return nil, err
	}

	return &Task{
		UserID:        userID,
		AssignmentID:  assignmentID,
		WeekID:        weekID,
		Title:         normalizedTitle,
		Description:   normalizedDescription,
		Status:        normalizedStatus,
		SpentHours:    spentHours,
		Observations:  normalizedObservations,
		WeekStartDate: normalizedWeekStartDate,
		Late:          late,
	}, nil
}

func (t *Task) UpdateTask(
	assignmentID uint,
	title, description string,
	status TaskStatus,
	spentHours int,
	observations string,
	weekStartDate time.Time,
	late bool,
) error {
	if t.Late {
		return ErrTaskLateUpdateForbidden
	}

	normalizedTitle := NormalizeTaskTitle(title)
	normalizedDescription := NormalizeTaskDescription(description)
	normalizedStatus := NormalizeTaskStatus(status)
	normalizedObservations := NormalizeTaskObservations(observations)
	normalizedWeekStartDate := NormalizeWeekStartDate(weekStartDate)

	if err := ValidateTaskAssignmentID(assignmentID); err != nil {
		return err
	}
	if err := ValidateTaskTitle(normalizedTitle); err != nil {
		return err
	}
	if err := ValidateTaskDescription(normalizedDescription); err != nil {
		return err
	}
	if err := ValidateTaskStatus(normalizedStatus); err != nil {
		return err
	}
	if err := ValidateTaskSpentHours(spentHours); err != nil {
		return err
	}
	if err := ValidateTaskWeekStartDate(normalizedWeekStartDate); err != nil {
		return err
	}

	t.AssignmentID = assignmentID
	t.Title = normalizedTitle
	t.Description = normalizedDescription
	t.Status = normalizedStatus
	t.SpentHours = spentHours
	t.Observations = normalizedObservations
	t.WeekStartDate = normalizedWeekStartDate
	t.Late = late

	return nil
}

func (t *Task) CanDelete(now time.Time) error {
	if !IsWeekActive(t.WeekStartDate, now) {
		return ErrTaskDeleteForbidden
	}
	return nil
}
