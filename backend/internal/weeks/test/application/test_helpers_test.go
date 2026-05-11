package application_test

import "backend/internal/weeks/domain"

type mockWeekRepository struct {
	weeks []domain.Week
}

func (m *mockWeekRepository) CreateMany(weeks []domain.Week) error {
	startID := uint(len(m.weeks) + 1)
	for i := range weeks {
		weeks[i].ID = startID + uint(i)
		m.weeks = append(m.weeks, weeks[i])
	}
	return nil
}

func (m *mockWeekRepository) FindAllByPeriodID(periodID uint) ([]domain.Week, error) {
	result := make([]domain.Week, 0)
	for _, week := range m.weeks {
		if week.PeriodID == periodID {
			result = append(result, week)
		}
	}
	return result, nil
}

func (m *mockWeekRepository) FindByPeriodIDAndNumber(periodID uint, number int) (*domain.Week, error) {
	for _, week := range m.weeks {
		if week.PeriodID == periodID && week.Number == number {
			copy := week
			return &copy, nil
		}
	}
	return nil, domain.ErrWeekNotFound
}

func (m *mockWeekRepository) FindByPeriodIDAndStartDate(periodID uint, startDate string) (*domain.Week, error) {
	for _, week := range m.weeks {
		if week.PeriodID == periodID && week.InitialDate == startDate {
			copy := week
			return &copy, nil
		}
	}
	return nil, domain.ErrWeekNotFound
}

func (m *mockWeekRepository) ExistsByPeriodID(periodID uint) (bool, error) {
	for _, week := range m.weeks {
		if week.PeriodID == periodID {
			return true, nil
		}
	}
	return false, nil
}
