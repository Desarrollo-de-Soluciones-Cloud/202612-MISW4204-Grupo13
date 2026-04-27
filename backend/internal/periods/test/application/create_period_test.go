package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	weeksApplication "backend/internal/weeks/application"
	"errors"
	"testing"
)

type MockPeriodRepository struct {
	periods    map[uint]*domain.Period
	periodsByName map[string]*domain.Period
	nextID     uint
}

type MockCreateWeeksForPeriod struct {
	lastInput *weeksApplication.CreateWeeksForPeriodInput
}

func (m *MockCreateWeeksForPeriod) Execute(input weeksApplication.CreateWeeksForPeriodInput) (*weeksApplication.CreateWeeksForPeriodOutput, error) {
	m.lastInput = &input
	return &weeksApplication.CreateWeeksForPeriodOutput{}, nil
}

func NewMockPeriodRepository() *MockPeriodRepository {
	return &MockPeriodRepository{
		periods:    make(map[uint]*domain.Period),
		periodsByName: make(map[string]*domain.Period),
		nextID:     1,
	}
}

func (m *MockPeriodRepository) Create(period *domain.Period) error {
	period.ID = m.nextID
	m.nextID++
	m.periods[period.ID] = period
	m.periodsByName[period.Name] = period
	return nil
}

func (m *MockPeriodRepository) FindByID(id uint) (*domain.Period, error) {
	if period, ok := m.periods[id]; ok {
		return period, nil
	}
	return nil, domain.ErrPeriodNotFound
}

func (m *MockPeriodRepository) FindByName(name string) (*domain.Period, error) {
	if period, ok := m.periodsByName[name]; ok {
		return period, nil
	}
	return nil, domain.ErrPeriodNotFound
}

func (m *MockPeriodRepository) FindAll() ([]domain.Period, error) {
	periods := make([]domain.Period, 0, len(m.periods))
	for _, p := range m.periods {
		periods = append(periods, *p)
	}
	return periods, nil
}

func (m *MockPeriodRepository) FindAllByState(state domain.PeriodState) ([]domain.Period, error) {
	periods := make([]domain.Period, 0)
	for _, p := range m.periods {
		if p.PeriodState == state {
			periods = append(periods, *p)
		}
	}
	return periods, nil
}

func (m *MockPeriodRepository) Update(period *domain.Period) error {
	if _, ok := m.periods[period.ID]; !ok {
		return domain.ErrPeriodNotFound
	}

	for name, existingPeriod := range m.periodsByName {
		if existingPeriod.ID == period.ID && name != period.Name {
			delete(m.periodsByName, name)
		}
	}

	m.periods[period.ID] = period
	m.periodsByName[period.Name] = period
	return nil
}

func (m *MockPeriodRepository) Delete(id uint) error {
	period, ok := m.periods[id]
	if !ok {
		return domain.ErrPeriodNotFound
	}
	delete(m.periods, id)
	delete(m.periodsByName, period.Name)
	return nil
}

func TestCreatePeriodSuccess(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	mockCreateWeeks := &MockCreateWeeksForPeriod{}
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, mockCreateWeeks)

	initialDate := "2026-10-05"
	weeksCount := 8

	input := applicationpkg.CreatePeriodInput{
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	}

	output, err := createPeriod.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2026-10" {
		t.Errorf("expected name '2026-10', got %q", output.Name)
	}
	if output.PeriodState != domain.ActivePeriod {
		t.Errorf("expected state %q, got %q", domain.ActivePeriod, output.PeriodState)
	}
	if output.WeeksCount != weeksCount {
		t.Errorf("expected weeks count %d, got %d", weeksCount, output.WeeksCount)
	}

	storedPeriod, err := mockRepo.FindByID(output.ID)
	if err != nil {
		t.Fatalf("expected stored period, got %v", err)
	}
	if storedPeriod.Name != "2026-10" {
		t.Errorf("expected stored period name, got %q", storedPeriod.Name)
	}
	if mockCreateWeeks.lastInput == nil {
		t.Fatal("expected weeks creation to be executed")
	}
	if mockCreateWeeks.lastInput.PeriodID != output.ID {
		t.Errorf("expected weeks period id %d, got %d", output.ID, mockCreateWeeks.lastInput.PeriodID)
	}
	if mockCreateWeeks.lastInput.WeeksCount != weeksCount {
		t.Errorf("expected weeks count %d, got %d", weeksCount, mockCreateWeeks.lastInput.WeeksCount)
	}
}

func TestCreatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, nil)

	initialDate := "2026-10-05"
	weeksCount := 8

	input := applicationpkg.CreatePeriodInput{
		Name:        "",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodNameRequired) {
		t.Errorf("expected ErrPeriodNameRequired, got %v", err)
	}
}

func TestCreatePeriodInvalidState(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, nil)

	initialDate := "2026-10-05"
	weeksCount := 8

	input := applicationpkg.CreatePeriodInput{
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.PeriodState("invalid"),
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestCreatePeriodNameAlreadyExists(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2026-10-05"
	weeksCount := 8

	existingPeriod, _ := domain.NewPeriod("2026-10", initialDate, weeksCount, domain.ActivePeriod)
	mockRepo.Create(existingPeriod)

	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, nil)

	input := applicationpkg.CreatePeriodInput{
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}

func TestCreatePeriodInvalidInitialDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, nil)

	weeksCount := 8

	input := applicationpkg.CreatePeriodInput{
		Name:        "2026-10",
		InitialDate: "2024",
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodInitialDateWrongFormat) {
		t.Errorf("expected ErrPeriodInitialDateWrongFormat, got %v", err)
	}
}

func TestCreatePeriodInvalidWeeksCount(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo, nil)

	initialDate := "2026-10-05"
	weeksCount := 10

	input := applicationpkg.CreatePeriodInput{
		Name:        "2026-10",
		InitialDate: initialDate,
		WeeksCount:  weeksCount,
		PeriodState: domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodWeeksCountInvalid) {
		t.Errorf("expected ErrPeriodWeeksCountInvalid, got %v", err)
	}
}
