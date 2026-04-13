package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestUpdatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	newInitialDate := "2026-10-12"
	newWeeksCount := 16

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	output, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period.ID,
		Name:        "2026-11",
		InitialDate: newInitialDate,
		WeeksCount:  newWeeksCount,
		PeriodState: domain.ClosedPeriod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2026-11" {
		t.Errorf("expected name '2026-11', got %q", output.Name)
	}
	if output.PeriodState != domain.ClosedPeriod {
		t.Errorf("expected state %q, got %q", domain.ClosedPeriod, output.PeriodState)
	}
	if output.WeeksCount != newWeeksCount {
		t.Errorf("expected weeks count %d, got %d", newWeeksCount, output.WeeksCount)
	}
}

func TestUpdatePeriodNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          999,
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestUpdatePeriodInvalidState(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period.ID,
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.PeriodState("invalid"),
	})

	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestUpdatePeriodNameAlreadyExists(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period1, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period1)

	period2, _ := domain.NewPeriod("2026-11", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period2)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period2.ID,
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}

func TestUpdatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period.ID,
		Name:        "invalid name",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNameWrongFormat) {
		t.Errorf("expected ErrPeriodNameWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidInitialDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period.ID,
		Name:        "2026-11",
		InitialDate: "2024",
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodInitialDateWrongFormat) {
		t.Errorf("expected ErrPeriodInitialDateWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidWeeksCount(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	period, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:          period.ID,
		Name:        "2026-11",
		InitialDate: initialDate,
		WeeksCount:  10,
		PeriodState: domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodWeeksCountInvalid) {
		t.Errorf("expected ErrPeriodWeeksCountInvalid, got %v", err)
	}
}
