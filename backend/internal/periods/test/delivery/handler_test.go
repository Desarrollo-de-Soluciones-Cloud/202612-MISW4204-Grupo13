package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/delivery"
	"backend/internal/periods/domain"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockPeriodRepository struct {
	periods map[uint]*domain.Period
	nextID  uint
}

func newMockPeriodRepository() *MockPeriodRepository {
	return &MockPeriodRepository{
		periods: make(map[uint]*domain.Period),
		nextID:  1,
	}
}

func (m *MockPeriodRepository) Create(period *domain.Period) error {
	period.ID = m.nextID
	m.periods[m.nextID] = period
	m.nextID++
	return nil
}

func (m *MockPeriodRepository) FindByID(id uint) (*domain.Period, error) {
	if p, exists := m.periods[id]; exists {
		return p, nil
	}
	return nil, domain.ErrPeriodNotFound
}

func (m *MockPeriodRepository) FindByName(name string) (*domain.Period, error) {
	for _, p := range m.periods {
		if p.Name == name {
			return p, nil
		}
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
	if _, exists := m.periods[period.ID]; !exists {
		return domain.ErrPeriodNotFound
	}
	m.periods[period.ID] = period
	return nil
}

func (m *MockPeriodRepository) Delete(id uint) error {
	if _, exists := m.periods[id]; !exists {
		return domain.ErrPeriodNotFound
	}
	delete(m.periods, id)
	return nil
}

func (m *MockPeriodRepository) AutoMigrate() error {
	return nil
}

func setupTestHandler(repo domain.PeriodRepository) *delivery.PeriodHandler {
	return delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
	)
}

func TestCreatePeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        "2026-10",
		InitialDate: "2026-10-05",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/periods", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusCreated {
		t.Errorf("CreatePeriod() status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp delivery.CreatePeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "2026-10" {
		t.Errorf("CreatePeriod() name = %s, want 2026-10", resp.Name)
	}
	if resp.WeeksCount != weeksCount {
		t.Errorf("CreatePeriod() weeks count = %d, want %d", resp.WeeksCount, weeksCount)
	}
}

func TestCreatePeriodBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        "", // Missing name
		InitialDate: "2026-10-05",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/periods", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreatePeriodInvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        "2026-10",
		InitialDate: "invalid-date",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/periods", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod() status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestListPeriodsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusOK {
		t.Errorf("ListPeriods() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.ListPeriodsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Periods) != 1 {
		t.Errorf("ListPeriods() returned %d periods, want 1", len(resp.Periods))
	}
	if resp.Periods[0].WeeksCount != 8 {
		t.Errorf("ListPeriods() weeks count = %d, want 8", resp.Periods[0].WeeksCount)
	}
}

func TestListPeriodsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusOK {
		t.Errorf("ListPeriods() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.ListPeriodsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Periods) != 0 {
		t.Errorf("ListPeriods() returned %d periods, want 0", len(resp.Periods))
	}
}

func TestGetPeriodByIDSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods/1", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetPeriodByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("GetPeriodByID() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "2026-10" {
		t.Errorf("GetPeriodByID() name = %s, want 2026-10", resp.Name)
	}
	if resp.WeeksCount != 8 {
		t.Errorf("GetPeriodByID() weeks count = %d, want 8", resp.WeeksCount)
	}
}

func TestGetPeriodByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods/999", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.GetPeriodByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetPeriodByID() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpdatePeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	repo.Create(period)

	weeksCount := 16
	body := delivery.UpdatePeriodRequest{
		Name:        "2026-10",
		InitialDate: "2026-10-05",
		WeeksCount:  &weeksCount,
		PeriodState: "closed",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/periods/1", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdatePeriod(c)

	if w.Code != http.StatusOK {
		t.Errorf("UpdatePeriod() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.PeriodState != "closed" {
		t.Errorf("UpdatePeriod() state = %s, want closed", resp.PeriodState)
	}
	if resp.WeeksCount != weeksCount {
		t.Errorf("UpdatePeriod() weeks count = %d, want %d", resp.WeeksCount, weeksCount)
	}
}

func TestUpdatePeriodNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.UpdatePeriodRequest{
		Name:        "2026-10",
		InitialDate: "2026-10-05",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/periods/999", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.UpdatePeriod(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("UpdatePeriod() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestDeletePeriodNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/periods/999", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.DeletePeriod(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("DeletePeriod() status = %d, want %d", w.Code, http.StatusNotFound)
	}
}