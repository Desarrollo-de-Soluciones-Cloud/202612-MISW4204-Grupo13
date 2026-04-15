package application

import (
	assignmentsDomain "backend/internal/assignments/domain"
	"backend/internal/tasks/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"errors"
	"time"
)

type taskContext struct {
	assignment    *assignmentsDomain.Assignment
	workspace     *workspacesDomain.Workspace
	week          *weeksDomain.Week
	weekStartDate time.Time
	weekFinalDate time.Time
}

func loadTaskContext(
	assignmentRepo TaskAssignmentRepository,
	workspaceRepo TaskWorkspaceRepository,
	weekRepo TaskWeekRepository,
	assignmentID uint,
	weekID uint,
) (*taskContext, error) {
	assignment, err := assignmentRepo.FindByID(assignmentID)
	if err != nil {
		if errors.Is(err, assignmentsDomain.ErrAssignmentNotFound) {
			return nil, domain.ErrTaskAssignmentNotFound
		}
		return nil, err
	}

	workspace, err := workspaceRepo.FindByID(assignment.WorkspaceID)
	if err != nil {
		if errors.Is(err, workspacesDomain.ErrWorkspaceNotFound) || err.Error() == workspacesDomain.ErrWorkspaceNotFound.Error() {
			return nil, domain.ErrTaskWorkspaceNotFound
		}
		return nil, err
	}

	week, err := weekRepo.FindByID(weekID)
	if err != nil {
		if errors.Is(err, weeksDomain.ErrWeekNotFound) {
			return nil, domain.ErrTaskWeekNotFound
		}
		return nil, err
	}

	if week.PeriodID != workspace.PeriodID {
		return nil, domain.ErrTaskWeekPeriodMismatch
	}

	weekStartDate, err := time.Parse("2006-01-02", week.InitialDate)
	if err != nil {
		return nil, domain.ErrTaskWeekStartDateInvalid
	}

	weekFinalDate, err := time.Parse("2006-01-02", week.FinalDate)
	if err != nil {
		return nil, domain.ErrTaskWeekStartDateInvalid
	}

	return &taskContext{
		assignment:    assignment,
		workspace:     workspace,
		week:          week,
		weekStartDate: weekStartDate,
		weekFinalDate: weekFinalDate,
	}, nil
}

func isClosedWeek(weekFinalDate, now time.Time) bool {
	return normalizeDateOnly(now).After(normalizeDateOnly(weekFinalDate))
}

func isActiveWeek(weekStartDate, weekFinalDate, now time.Time) bool {
	currentDate := normalizeDateOnly(now)
	return (currentDate.Equal(normalizeDateOnly(weekStartDate)) || currentDate.After(normalizeDateOnly(weekStartDate))) &&
		(currentDate.Equal(normalizeDateOnly(weekFinalDate)) || currentDate.Before(normalizeDateOnly(weekFinalDate)))
}

func normalizeDateOnly(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}
