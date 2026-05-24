package infrastructure_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	sharedDB "backend/internal/shared/database/testsupport"
	tasksDomain "backend/internal/tasks/domain"
	"testing"
	"time"
)

func TestTaskReaderFindAllByWorkspaceAndWeekSQLite(t *testing.T) {
	db := sharedDB.SetupSQLiteDB(t, &assignmentsDomain.Assignment{}, &tasksDomain.Task{})

	assignments := []assignmentsDomain.Assignment{
		{UserID: 10, WorkspaceID: 1, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 6},
		{UserID: 11, WorkspaceID: 2, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 6},
	}
	for i := range assignments {
		if err := db.Create(&assignments[i]).Error; err != nil {
			t.Fatalf("expected assignment insert, got %v", err)
		}
	}

	tasks := []tasksDomain.Task{
		{AssignmentID: assignments[0].ID, Title: "A", Description: "A", Status: tasksDomain.TaskStatusEnDesarrollo, SpentHours: 2, WeekID: uintPtr(7), WeekStartDate: mustTaskDate(t, "2026-10-05")},
		{AssignmentID: assignments[0].ID, Title: "B", Description: "B", Status: tasksDomain.TaskStatusFinalizado, SpentHours: 3, WeekStartDate: mustTaskDate(t, "2026-10-05")},
		{AssignmentID: assignments[1].ID, Title: "C", Description: "C", Status: tasksDomain.TaskStatusEnDesarrollo, SpentHours: 1, WeekID: uintPtr(7), WeekStartDate: mustTaskDate(t, "2026-10-05")},
	}
	for i := range tasks {
		if err := db.Create(&tasks[i]).Error; err != nil {
			t.Fatalf("expected task insert, got %v", err)
		}
	}

	reader := reportsInfrastructure.NewTaskReader()
	found, err := reader.FindAllByWorkspaceAndWeek(1, 7, "2026-10-05")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(found) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(found))
	}
}

func mustTaskDate(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("expected valid task date, got %v", err)
	}
	return parsed
}

func uintPtr(value uint) *uint {
	return &value
}
