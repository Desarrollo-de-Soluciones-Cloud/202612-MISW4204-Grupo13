package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
	"time"
)

func TestUpdatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	newInitialDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	newFinalDate := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	newInscriptionDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	output, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-02",
		InitialDate:          newInitialDate,
		FinalDate:            newFinalDate,
		InscriptionFinalDate: newInscriptionDate,
		PeriodState:          domain.ClosedPeriod,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2024-02" {
		t.Errorf("expected name '2024-02', got %q", output.Name)
	}
	if output.PeriodState != domain.ClosedPeriod {
		t.Errorf("expected state %q, got %q", domain.ClosedPeriod, output.PeriodState)
	}
}

func TestUpdatePeriodNotFound(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   999,
		Name:                 "2024-02",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNotFound) {
		t.Errorf("expected ErrPeriodNotFound, got %v", err)
	}
}

func TestUpdatePeriodInvalidState(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	finalDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	inscriptionDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.PeriodState("invalid"),
	})

	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}
