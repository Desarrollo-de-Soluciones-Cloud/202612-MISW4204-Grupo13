package domain

import (
	domainpkg "backend/internal/tasks/domain"
	"errors"
	"testing"
	"time"
)

func TestNewTaskSuccess(t *testing.T) {
	weekDate := time.Date(2026, 4, 8, 10, 0, 0, 0, time.UTC)

	task, err := domainpkg.NewTask(
		1,
		10,
		nil,
		" Prepare class ",
		" Review slides ",
		domainpkg.TaskStatusAbierto,
		2,
		" bring examples ",
		weekDate,
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedMonday := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
	if task.AssignmentID != 10 {
		t.Fatalf("expected assignment id 10, got %d", task.AssignmentID)
	}
	if task.Title != "Prepare class" {
		t.Fatalf("expected normalized title, got %q", task.Title)
	}
	if task.Description != "Review slides" {
		t.Fatalf("expected normalized description, got %q", task.Description)
	}
	if task.Observations != "bring examples" {
		t.Fatalf("expected normalized observations, got %q", task.Observations)
	}
	if !task.WeekStartDate.Equal(expectedMonday) {
		t.Fatalf("expected monday %v, got %v", expectedMonday, task.WeekStartDate)
	}
}

func TestNewTaskRejectsMissingAssignmentID(t *testing.T) {
	_, err := domainpkg.NewTask(
		1,
		0,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		false,
	)
	if !errors.Is(err, domainpkg.ErrTaskAssignmentIDRequired) {
		t.Fatalf("expected ErrTaskAssignmentIDRequired, got %v", err)
	}
}

func TestNewTaskAcceptsOfficialStatuses(t *testing.T) {
	validStatuses := []domainpkg.TaskStatus{
		domainpkg.TaskStatusAbierto,
		domainpkg.TaskStatusEnDesarrollo,
		domainpkg.TaskStatusFinalizado,
	}

	for _, status := range validStatuses {
		task, err := domainpkg.NewTask(
			1,
			10,
			nil,
			"Prepare class",
			"Review slides",
			status,
			2,
			"",
			time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
			false,
		)
		if err != nil {
			t.Fatalf("expected status %q to be valid, got %v", status, err)
		}
		if task.Status != status {
			t.Fatalf("expected task status %q, got %q", status, task.Status)
		}
	}
}

func TestNewTaskRejectsLegacyStatus(t *testing.T) {
	_, err := domainpkg.NewTask(
		1,
		10,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatus("open"),
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		false,
	)
	if !errors.Is(err, domainpkg.ErrTaskStatusInvalid) {
		t.Fatalf("expected ErrTaskStatusInvalid, got %v", err)
	}
}

func TestNewTaskRejectsSpentHoursBelowOne(t *testing.T) {
	_, err := domainpkg.NewTask(
		1,
		10,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		-1,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		false,
	)
	if !errors.Is(err, domainpkg.ErrTaskSpentHoursInvalid) {
		t.Fatalf("expected ErrTaskSpentHoursInvalid, got %v", err)
	}
}

func TestNormalizeWeekStartDateMovesAnyDayToMonday(t *testing.T) {
	date := time.Date(2026, 4, 12, 15, 0, 0, 0, time.UTC)

	normalized := domainpkg.NormalizeWeekStartDate(date)
	expected := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)

	if !normalized.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, normalized)
	}
}

func TestIsWeekClosed(t *testing.T) {
	weekStartDate := time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC)
	now := time.Date(2026, 4, 13, 9, 0, 0, 0, time.UTC)

	if !domainpkg.IsWeekClosed(weekStartDate, now) {
		t.Fatal("expected week to be closed")
	}
}

func TestUpdateTaskRejectsLateTask(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		true,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.UpdateTask(
		10,
		"Prepare class 2",
		"Review slides",
		domainpkg.TaskStatusFinalizado,
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		true,
	)
	if !errors.Is(err, domainpkg.ErrTaskLateUpdateForbidden) {
		t.Fatalf("expected ErrTaskLateUpdateForbidden, got %v", err)
	}
}

func TestCanDeleteRejectsInactiveWeek(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.CanDelete(time.Date(2026, 4, 13, 9, 0, 0, 0, time.UTC))
	if !errors.Is(err, domainpkg.ErrTaskDeleteForbidden) {
		t.Fatalf("expected ErrTaskDeleteForbidden, got %v", err)
	}
}

func TestUpdateTaskAllowsOfficialStatusChange(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		nil,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.UpdateTask(
		10,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusEnDesarrollo,
		2,
		"",
		time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected status update to succeed, got %v", err)
	}
	if task.Status != domainpkg.TaskStatusEnDesarrollo {
		t.Fatalf("expected updated status %q, got %q", domainpkg.TaskStatusEnDesarrollo, task.Status)
	}
}
