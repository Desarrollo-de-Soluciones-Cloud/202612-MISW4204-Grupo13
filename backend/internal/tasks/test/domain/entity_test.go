package domain

import (
	domainpkg "backend/internal/tasks/domain"
	"errors"
	"testing"
	"time"
)

func TestNewTaskSuccess(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		30,
		" Prepare class ",
		" Review slides ",
		domainpkg.TaskStatusAbierto,
		2,
		" bring examples ",
		[]domainpkg.TaskAttachment{{Path: " docs/guide.pdf "}},
		time.Date(2026, 4, 13, 10, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedWeekStart := time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC)
	if task.AssignmentID != 10 {
		t.Fatalf("expected assignment id 10, got %d", task.AssignmentID)
	}
	if task.WeekID != 30 {
		t.Fatalf("expected week id 30, got %d", task.WeekID)
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
	if !task.WeekStartDate.Equal(expectedWeekStart) {
		t.Fatalf("expected week start %v, got %v", expectedWeekStart, task.WeekStartDate)
	}
	if len(task.Attachments) != 1 || task.Attachments[0].Path != "docs/guide.pdf" {
		t.Fatalf("expected normalized attachment, got %+v", task.Attachments)
	}
}

func TestNewTaskRejectsMissingWeekID(t *testing.T) {
	_, err := domainpkg.NewTask(
		1,
		10,
		0,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if !errors.Is(err, domainpkg.ErrTaskWeekIDRequired) {
		t.Fatalf("expected ErrTaskWeekIDRequired, got %v", err)
	}
}

func TestNewTaskRejectsAttachmentWithoutPath(t *testing.T) {
	_, err := domainpkg.NewTask(
		1,
		10,
		30,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		[]domainpkg.TaskAttachment{{Path: " "}},
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if !errors.Is(err, domainpkg.ErrTaskAttachmentPathRequired) {
		t.Fatalf("expected ErrTaskAttachmentPathRequired, got %v", err)
	}
}

func TestUpdateTaskRejectsLateTask(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		30,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		true,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.UpdateTask(
		10,
		30,
		"Prepare class 2",
		"Review slides",
		domainpkg.TaskStatusFinalizado,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		true,
	)
	if !errors.Is(err, domainpkg.ErrTaskLateUpdateForbidden) {
		t.Fatalf("expected ErrTaskLateUpdateForbidden, got %v", err)
	}
}

func TestUpdateTaskSuccess(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		30,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.UpdateTask(
		10,
		30,
		"Prepare class 2",
		"Review slides updated",
		domainpkg.TaskStatusFinalizado,
		3,
		"done",
		[]domainpkg.TaskAttachment{{Path: "docs/result.pdf"}},
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if task.Title != "Prepare class 2" || task.Status != domainpkg.TaskStatusFinalizado {
		t.Fatalf("expected task to be updated, got %+v", task)
	}
	if len(task.Attachments) != 1 || task.Attachments[0].Path != "docs/result.pdf" {
		t.Fatalf("expected attachments to be updated, got %+v", task.Attachments)
	}
}

func TestCanDeleteRejectsInactiveWeek(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		30,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = task.CanDelete(false)
	if !errors.Is(err, domainpkg.ErrTaskDeleteForbidden) {
		t.Fatalf("expected ErrTaskDeleteForbidden, got %v", err)
	}
}

func TestCanDeleteAllowsActiveWeek(t *testing.T) {
	task, err := domainpkg.NewTask(
		1,
		10,
		30,
		"Prepare class",
		"Review slides",
		domainpkg.TaskStatusAbierto,
		2,
		"",
		nil,
		time.Date(2026, 4, 13, 0, 0, 0, 0, time.UTC),
		false,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := task.CanDelete(true); err != nil {
		t.Fatalf("expected delete to be allowed, got %v", err)
	}
}

func TestValidateTaskStatusAndHours(t *testing.T) {
	if err := domainpkg.ValidateTaskStatus(domainpkg.TaskStatus("")); !errors.Is(err, domainpkg.ErrTaskStatusRequired) {
		t.Fatalf("expected ErrTaskStatusRequired, got %v", err)
	}
	if err := domainpkg.ValidateTaskStatus(domainpkg.TaskStatus("open")); !errors.Is(err, domainpkg.ErrTaskStatusInvalid) {
		t.Fatalf("expected ErrTaskStatusInvalid, got %v", err)
	}
	if err := domainpkg.ValidateTaskSpentHours(0); !errors.Is(err, domainpkg.ErrTaskSpentHoursRequired) {
		t.Fatalf("expected ErrTaskSpentHoursRequired, got %v", err)
	}
	if err := domainpkg.ValidateTaskSpentHours(-1); !errors.Is(err, domainpkg.ErrTaskSpentHoursInvalid) {
		t.Fatalf("expected ErrTaskSpentHoursInvalid, got %v", err)
	}
	if !domainpkg.IsValidTaskStatus(domainpkg.TaskStatusEnDesarrollo) {
		t.Fatal("expected en desarrollo to be valid")
	}
}
