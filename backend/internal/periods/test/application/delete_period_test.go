package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
	"time"
)

func TestDeletePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
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
