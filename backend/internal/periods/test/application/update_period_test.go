package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

func TestUpdatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	newInitialDate := "2024-02-01"
	newFinalDate := "2024-07-31"
	newInscriptionDate := "2024-02-15"

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

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

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

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

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

func TestUpdatePeriodNameAlreadyExists(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period1, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period1)

	period2, _ := domain.NewPeriod("2024-02", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period2)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period2.ID,
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}

func TestUpdatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "invalid name",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodNameWrongFormat) {
		t.Errorf("expected ErrPeriodNameWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidInitialDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-02",
		InitialDate:          "2024",
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodInitialDateWrongFormat) {
		t.Errorf("expected ErrPeriodInitialDateWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidFinalDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-02",
		InitialDate:          initialDate,
		FinalDate:            "2024",
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodFinalDateWrongFormat) {
		t.Errorf("expected ErrPeriodFinalDateWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidInscriptionDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-02",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: "2024",
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodInscriptionFinalDateWrongFormat) {
		t.Errorf("expected ErrPeriodInscriptionFinalDateWrongFormat, got %v", err)
	}
}

func TestUpdatePeriodInvalidDateSequence(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	period, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(period)

	updatePeriod := applicationpkg.NewUpdatePeriod(mockRepo)
	_, err := updatePeriod.Execute(applicationpkg.UpdatePeriodInput{
		ID:                   period.ID,
		Name:                 "2024-02",
		InitialDate:          "2024-06-30",
		FinalDate:            "2024-01-01",
		InscriptionFinalDate: "2024-01-15",
		PeriodState:          domain.ActivePeriod,
	})

	if !errors.Is(err, domain.ErrPeriodDateSequenceInvalid) {
		t.Errorf("expected ErrPeriodDateSequenceInvalid, got %v", err)
	}
}
