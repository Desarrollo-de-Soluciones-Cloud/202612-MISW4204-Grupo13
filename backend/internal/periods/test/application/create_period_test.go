package application

import (
	applicationpkg "backend/internal/periods/application"
	"backend/internal/periods/domain"
	"errors"
	"testing"
)

type MockPeriodRepository struct {
	periods    map[uint]*domain.Period
	periodsByName map[string]*domain.Period
	nextID     uint
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
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	output, err := createPeriod.Execute(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.Name != "2024-01" {
		t.Errorf("expected name '2024-01', got %q", output.Name)
	}
	if output.PeriodState != domain.ActivePeriod {
		t.Errorf("expected state %q, got %q", domain.ActivePeriod, output.PeriodState)
	}

	storedPeriod, err := mockRepo.FindByID(output.ID)
	if err != nil {
		t.Fatalf("expected stored period, got %v", err)
	}
	if storedPeriod.Name != "2024-01" {
		t.Errorf("expected stored period name, got %q", storedPeriod.Name)
	}
}

func TestCreatePeriodInvalidName(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "ab",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodNameWrongFormat) {
		t.Errorf("expected ErrPeriodNameWrongFormat, got %v", err)
	}
}

func TestCreatePeriodInvalidDateSequence(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-06-30"
	finalDate := "2024-01-01"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodDateSequenceInvalid) {
		t.Errorf("expected ErrPeriodDateSequenceInvalid, got %v", err)
	}
}

func TestCreatePeriodInvalidState(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.PeriodState("invalid"),
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodStateInvalid) {
		t.Errorf("expected ErrPeriodStateInvalid, got %v", err)
	}
}

func TestCreatePeriodNameAlreadyExists(t *testing.T) {
	mockRepo := NewMockPeriodRepository()

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	existingPeriod, _ := domain.NewPeriod("2024-01", initialDate, finalDate, inscriptionDate, domain.ActivePeriod)
	mockRepo.Create(existingPeriod)

	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodNameAlreadyExists) {
		t.Errorf("expected ErrPeriodNameAlreadyExists, got %v", err)
	}
}

func TestCreatePeriodInvalidInitialDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	finalDate := "2024-06-30"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          "2024",
		FinalDate:            finalDate,
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodInitialDateWrongFormat) {
		t.Errorf("expected ErrPeriodInitialDateWrongFormat, got %v", err)
	}
}

func TestCreatePeriodInvalidFinalDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-01-01"
	inscriptionDate := "2024-01-15"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            "2024",
		InscriptionFinalDate: inscriptionDate,
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodFinalDateWrongFormat) {
		t.Errorf("expected ErrPeriodFinalDateWrongFormat, got %v", err)
	}
}

func TestCreatePeriodInvalidInscriptionDateFormat(t *testing.T) {
	mockRepo := NewMockPeriodRepository()
	createPeriod := applicationpkg.NewCreatePeriod(mockRepo)

	initialDate := "2024-01-01"
	finalDate := "2024-06-30"

	input := applicationpkg.CreatePeriodInput{
		Name:                 "2024-01",
		InitialDate:          initialDate,
		FinalDate:            finalDate,
		InscriptionFinalDate: "2024",
		PeriodState:          domain.ActivePeriod,
	}

	_, err := createPeriod.Execute(input)
	if !errors.Is(err, domain.ErrPeriodInscriptionFinalDateWrongFormat) {
		t.Errorf("expected ErrPeriodInscriptionFinalDateWrongFormat, got %v", err)
	}
}
