package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestDeletePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

		initialDate := "2026-10-05"

		period, _ := domain.NewPeriod("2026-10", initialDate, 16, domain.ActivePeriod)
	mockRepo.Create(period)

	deletePeriod := applicationpkg.NewDeletePeriod(mockRepo)
	err := deletePeriod.Execute(applicationpkg.DeletePeriodInput{ID: period.ID})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = mockRepo.FindByID(period.ID)
	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected period to be deleted, but it still exists")
	}
}

func TestDeletePeriodNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	deletePeriod := applicationpkg.NewDeletePeriod(mockRepo)

	err := deletePeriod.Execute(applicationpkg.DeletePeriodInput{ID: 999})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}
