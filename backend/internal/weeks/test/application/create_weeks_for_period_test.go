package application_test

import (
	application "backend/internal/weeks/application"
	"backend/internal/weeks/domain"
	"errors"
	"testing"
)

type MockWeekRepository struct {
	weeks []domain.Week
}

func (m *MockWeekRepository) CreateMany(weeks []domain.Week) error {
	startID := uint(len(m.weeks) + 1)
	for i := range weeks {
		weeks[i].ID = startID + uint(i)
		m.weeks = append(m.weeks, weeks[i])
	}
	return nil
}

func (m *MockWeekRepository) FindAllByPeriodID(periodID uint) ([]domain.Week, error) {
	result := make([]domain.Week, 0)
	for _, week := range m.weeks {
		if week.PeriodID == periodID {
			result = append(result, week)
		}
	}
	return result, nil
}

func (m *MockWeekRepository) FindByPeriodIDAndNumber(periodID uint, number int) (*domain.Week, error) {
	for _, week := range m.weeks {
		if week.PeriodID == periodID && week.Number == number {
			copy := week
			return &copy, nil
		}
	}
	return nil, domain.ErrWeekNotFound
}

func (m *MockWeekRepository) ExistsByPeriodID(periodID uint) (bool, error) {
	for _, week := range m.weeks {
		if week.PeriodID == periodID {
			return true, nil
		}
	}
	return false, nil
}

func TestCreateWeeksForPeriodSuccess(t *testing.T) {
	repo := &MockWeekRepository{}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	output, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-03-08",
		WeeksCount:  8,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(output.Weeks) != 8 {
		t.Fatalf("expected 8 weeks, got %d", len(output.Weeks))
	}
	if output.Weeks[0].Number != 1 || output.Weeks[0].InitialDate != "2026-01-12" || output.Weeks[0].FinalDate != "2026-01-18" {
		t.Fatalf("unexpected first week: %+v", output.Weeks[0])
	}
	if output.Weeks[7].Number != 8 {
		t.Fatalf("expected last week number 8, got %d", output.Weeks[7].Number)
	}
}

func TestCreateWeeksForPeriodRejectsDuplicateGeneration(t *testing.T) {
	repo := &MockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
		},
	}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	_, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-03-08",
		WeeksCount:  8,
	})
	if !errors.Is(err, domain.ErrWeeksAlreadyExistForPeriod) {
		t.Fatalf("expected ErrWeeksAlreadyExistForPeriod, got %v", err)
	}
}

func TestCreateWeeksForPeriodRejectsInvalidInitialDate(t *testing.T) {
	repo := &MockWeekRepository{}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	_, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-13",
		FinalDate:   "2026-03-08",
		WeeksCount:  8,
	})
	if !errors.Is(err, domain.ErrWeekInitialDateMustBeMonday) {
		t.Fatalf("expected ErrWeekInitialDateMustBeMonday, got %v", err)
	}
}

func TestCreateWeeksForPeriodRejectsInvalidWeeksCount(t *testing.T) {
	repo := &MockWeekRepository{}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	_, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-03-08",
		WeeksCount:  12,
	})
	if !errors.Is(err, domain.ErrWeekCountInvalid) {
		t.Fatalf("expected ErrWeekCountInvalid, got %v", err)
	}
}

func TestCreateWeeksForPeriodRejectsFinalDateMismatch(t *testing.T) {
	repo := &MockWeekRepository{}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	_, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-03-15",
		WeeksCount:  8,
	})
	if !errors.Is(err, domain.ErrWeekFinalDateMismatch) {
		t.Fatalf("expected ErrWeekFinalDateMismatch, got %v", err)
	}
}

func TestCreateWeeksForPeriodRejectsFinalDateThatIsNotSunday(t *testing.T) {
	repo := &MockWeekRepository{}
	createWeeks := application.NewCreateWeeksForPeriod(repo)

	_, err := createWeeks.Execute(application.CreateWeeksForPeriodInput{
		PeriodID:    1,
		InitialDate: "2026-01-12",
		FinalDate:   "2026-03-09",
		WeeksCount:  8,
	})
	if !errors.Is(err, domain.ErrWeekFinalDateMustBeSunday) {
		t.Fatalf("expected ErrWeekFinalDateMustBeSunday, got %v", err)
	}
}
