package delivery

import (
	"backend/internal/periods/application"
	"backend/internal/periods/delivery"
	"backend/internal/periods/domain"
	weeksApplication "backend/internal/weeks/application"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockPeriodRepository struct {
	periods       map[uint]*domain.Period
	periodsByName map[string]*domain.Period
	nextID        uint
}

type MockCreateWeeksForPeriod struct{}

func (m *MockCreateWeeksForPeriod) Execute(input weeksApplication.CreateWeeksForPeriodInput) (*weeksApplication.CreateWeeksForPeriodOutput, error) {
	return &weeksApplication.CreateWeeksForPeriodOutput{}, nil
}

func newMockPeriodRepository() *MockPeriodRepository {
	return &MockPeriodRepository{
		periods:       make(map[uint]*domain.Period),
		periodsByName: make(map[string]*domain.Period),
		nextID:        1,
	}
}

func (m *MockPeriodRepository) Create(period *domain.Period) error {
	period.ID = m.nextID
	m.periods[m.nextID] = period
	m.periodsByName[period.Name] = period
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
	if period, ok := m.periodsByName[name]; ok {
		return period, nil
	}
	return nil, domain.ErrPeriodNotFound
}

func (m *MockPeriodRepository) FindAll() ([]domain.Period, error) {
	return collectPeriods(m.periods, nil), nil
}

func (m *MockPeriodRepository) FindAllByState(state domain.PeriodState) ([]domain.Period, error) {
	return collectPeriods(m.periods, func(period *domain.Period) bool {
		return period.PeriodState == state
	}), nil
}

func (m *MockPeriodRepository) Update(period *domain.Period) error {
	if _, exists := m.periods[period.ID]; !exists {
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
	period, exists := m.periods[id]
	if !exists {
		return domain.ErrPeriodNotFound
	}
	delete(m.periods, id)
	delete(m.periodsByName, period.Name)
	return nil
}

func (m *MockPeriodRepository) AutoMigrate() error {
	return nil
}

func setupTestHandler(repo domain.PeriodRepository) *delivery.PeriodHandler {
	createWeeks := &MockCreateWeeksForPeriod{}
	return delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo, createWeeks),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)
}

func collectPeriods(
	periods map[uint]*domain.Period,
	keep func(*domain.Period) bool,
) []domain.Period {
	result := make([]domain.Period, 0, len(periods))
	for _, period := range periods {
		if keep != nil && !keep(period) {
			continue
		}
		result = append(result, *period)
	}
	return result
}

func TestCreatePeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        testPeriod202610Name,
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusCreated {
		t.Errorf(errCreatePeriodStatus, w.Code, http.StatusCreated)
	}

	var resp delivery.CreatePeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != testPeriod202610Name {
		t.Errorf("CreatePeriod() name = %s, want %s", resp.Name, testPeriod202610Name)
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
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

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
		Name:        testPeriod202610Name,
		InitialDate: "invalid-date",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

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

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", testPeriodsPath, nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusOK {
		t.Errorf(errListPeriodsStatus, w.Code, http.StatusOK)
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
	req, _ := http.NewRequest("GET", testPeriodsPath, nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusOK {
		t.Errorf(errListPeriodsStatus, w.Code, http.StatusOK)
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

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", testPeriodByIDPath, nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetPeriodByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("GetPeriodByID() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != testPeriod202610Name {
		t.Errorf("GetPeriodByID() name = %s, want %s", resp.Name, testPeriod202610Name)
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

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	body := delivery.UpdatePeriodRequest{
		Name: testPeriod202611Name,
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", testPeriodByIDPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdatePeriod(c)

	if w.Code != http.StatusOK {
		t.Errorf("UpdatePeriod() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != testPeriod202611Name {
		t.Errorf("UpdatePeriod() name = %s, want %s", resp.Name, testPeriod202611Name)
	}
	// State and other fields should remain unchanged
	if resp.PeriodState != "active" {
		t.Errorf("UpdatePeriod() state = %s, want active", resp.PeriodState)
	}
}

func TestUpdatePeriodNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	body := delivery.UpdatePeriodRequest{
		Name: testPeriod202610Name,
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/periods/999", bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.UpdatePeriod(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("UpdatePeriod(not found) status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestListPeriodsByState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period1, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	period2, _ := domain.NewPeriod(testPeriod202611Name, testPeriodInitialDate1012, 16, domain.ActivePeriod)
	period3, _ := domain.NewPeriod("2026-12", "2026-11-09", 8, domain.ClosedPeriod)
	repo.Create(period1)
	repo.Create(period2)
	repo.Create(period3)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", testPeriodsPath+"?state=active", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusOK {
		t.Errorf("ListPeriods() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.ListPeriodsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Periods) != 2 {
		t.Errorf("ListPeriods(state=active) returned %d periods, want 2", len(resp.Periods))
	}
	for _, p := range resp.Periods {
		if p.PeriodState != "active" {
			t.Errorf("ListPeriods(state=active) returned period with state %s", p.PeriodState)
		}
	}
}

// TestListPeriodsByStateInvalid tests with invalid state parameter
func TestListPeriodsByStateInvalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", testPeriodsPath+"?state=invalidstate", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListPeriods(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ListPeriods(state=invalidstate) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestGetPeriodByIDInvalidFormat tests with invalid ID format
func TestGetPeriodByIDInvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods/invalid", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "invalid"})

	handler.GetPeriodByID(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("GetPeriodByID(invalid) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestCreatePeriodWithNilWeeksCount tests CreatePeriod with nil WeeksCount
func TestCreatePeriodWithNilWeeksCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	body := delivery.CreatePeriodRequest{
		Name:        testPeriod202610Name,
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  nil, // Missing WeeksCount
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod(nil WeeksCount) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreatePeriodInvalidName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        "",
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod(invalid name) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestCreatePeriodInvalidState tests CreatePeriod with invalid state
func TestCreatePeriodInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 8
	body := delivery.CreatePeriodRequest{
		Name:        testPeriod202610Name,
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  &weeksCount,
		PeriodState: "invalidstate",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod(invalid state) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestCreatePeriodInvalidWeeksCount tests CreatePeriod with invalid weeks count
func TestCreatePeriodInvalidWeeksCount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	weeksCount := 10 // Only 8 or 16 allowed
	body := delivery.CreatePeriodRequest{
		Name:        testPeriod202610Name,
		InitialDate: testPeriodInitialDate1005,
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", testPeriodsPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("CreatePeriod(invalid weeks count) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestUpdatePeriodInvalidID tests UpdatePeriod with invalid ID format
func TestUpdatePeriodInvalidState(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	body := delivery.UpdatePeriodRequest{
		Name: "",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", testPeriodByIDPath, bytes.NewBuffer(bodyJSON))
	req.Header.Set(testHeaderContentType, testJSONContentType)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdatePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("UpdatePeriod(invalid name) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestClosePeriodSuccess tests ClosePeriod successfully closing an active period
func TestClosePeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ActivePeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", testPeriodByIDPath+"/close", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.ClosePeriod(c)

	if w.Code != http.StatusOK {
		t.Errorf("ClosePeriod() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.PeriodState != "closed" {
		t.Errorf("ClosePeriod() state = %s, want closed", resp.PeriodState)
	}
}

// TestClosePeriodNotFound tests ClosePeriod with non-existent period
func TestClosePeriodNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/periods/999/close", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.ClosePeriod(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("ClosePeriod(not found) status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// TestClosePeriodAlreadyClosed tests ClosePeriod with already closed period
func TestClosePeriodAlreadyClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	period, _ := domain.NewPeriod(testPeriod202610Name, testPeriodInitialDate1005, 8, domain.ClosedPeriod)
	repo.Create(period)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", testPeriodByIDPath+"/close", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.ClosePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ClosePeriod(already closed) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// TestClosePeriodInvalidID tests ClosePeriod with invalid ID format
func TestClosePeriodInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockPeriodRepository()
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/periods/invalid/close", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "invalid"})

	handler.ClosePeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ClosePeriod(invalid id) status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
