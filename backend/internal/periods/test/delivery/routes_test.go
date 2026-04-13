package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/periods/application"
	"backend/internal/periods/delivery"
	"backend/internal/periods/domain"

	"github.com/gin-gonic/gin"
)

// TestSetupRoutesRegistersHandlers verifies that SetupRoutes registers all handlers correctly
func TestSetupRoutesRegistersHandlers(t *testing.T) {
	t.Skip("Skipping because it requires database setup")
	// This test verifies that routes are properly registered
	// In a real scenario, we would use a test database
}

// TestRoutesIntegration performs integration tests for all route handlers
func TestRoutesIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup mock repository
	repo := newMockPeriodRepository()

	// Create handlers
	createPeriod := application.NewCreatePeriod(repo)
	listPeriods := application.NewListPeriods(repo)
	listPeriodsByState := application.NewListPeriodsByState(repo)
	getPeriodByID := application.NewGetPeriodByID(repo)
	updatePeriod := application.NewUpdatePeriod(repo)

	handler := delivery.NewPeriodHandler(createPeriod, listPeriods, listPeriodsByState, getPeriodByID, updatePeriod, application.NewClosePeriod(repo))

	// Register routes manually (since SetupRoutes requires database)
	periods := router.Group("/periods")
	{
		periods.POST("", handler.CreatePeriod)
		periods.GET("", handler.ListPeriods)
		periods.GET("/:id", handler.GetPeriodByID)
		periods.PUT("/:id", handler.UpdatePeriod)
	}

	// Test POST /periods (Create)
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

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("POST /periods status = %d, want %d", w.Code, http.StatusCreated)
	}

	// Test GET /periods (List)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/periods", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /periods status = %d, want %d", w.Code, http.StatusOK)
	}

	// Test GET /periods/1 (Get by ID)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/periods/1", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /periods/1 status = %d, want %d", w.Code, http.StatusOK)
	}
}

// TestRoutesPostCreatePeriod tests the POST route for creating periods
func TestRoutesPostCreatePeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := newMockPeriodRepository()
	handler := delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)

	periods := router.Group("/periods")
	periods.POST("", handler.CreatePeriod)

	weeksCount := 16
	body := delivery.CreatePeriodRequest{
		Name:        "2026-11",
		InitialDate: "2026-10-12",
		WeeksCount:  &weeksCount,
		PeriodState: "active",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/periods", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("POST /periods status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp delivery.CreatePeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "2026-11" || resp.WeeksCount != 16 {
		t.Errorf("POST /periods response incorrect")
	}
}

// TestRoutesGetPeriods tests the GET route for listing periods
func TestRoutesGetPeriods(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := newMockPeriodRepository()
	
	// Add some test data
	period1, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2026-11", "2026-10-12", 16, domain.ClosedPeriod)
	repo.Create(period1)
	repo.Create(period2)

	handler := delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)

	periods := router.Group("/periods")
	periods.GET("", handler.ListPeriods)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /periods status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.ListPeriodsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Periods) != 2 {
		t.Errorf("GET /periods returned %d periods, want 2", len(resp.Periods))
	}
}

// TestRoutesGetPeriodsWithStateFilter tests the GET route with state filter
func TestRoutesGetPeriodsWithStateFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := newMockPeriodRepository()
	
	// Add test data
	period1, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	period2, _ := domain.NewPeriod("2026-11", "2026-10-12", 16, domain.ClosedPeriod)
	repo.Create(period1)
	repo.Create(period2)

	handler := delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)

	periods := router.Group("/periods")
	periods.GET("", handler.ListPeriods)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods?state=active", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /periods?state=active status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.ListPeriodsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Periods) != 1 {
		t.Errorf("GET /periods?state=active returned %d periods, want 1", len(resp.Periods))
	}
	if resp.Periods[0].PeriodState != "active" {
		t.Errorf("GET /periods?state=active returned non-active period")
	}
}

// TestRoutesGetPeriodByID tests the GET route for fetching a specific period
func TestRoutesGetPeriodByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := newMockPeriodRepository()
	
	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	repo.Create(period)

	handler := delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)

	periods := router.Group("/periods")
	periods.GET("/:id", handler.GetPeriodByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/periods/1", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /periods/1 status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "2026-10" {
		t.Errorf("GET /periods/1 returned wrong period")
	}
}

// TestRoutesPutUpdatePeriod tests the PUT route for updating periods
func TestRoutesPutUpdatePeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := newMockPeriodRepository()
	
	period, _ := domain.NewPeriod("2026-10", "2026-10-05", 8, domain.ActivePeriod)
	repo.Create(period)

	handler := delivery.NewPeriodHandler(
		application.NewCreatePeriod(repo),
		application.NewListPeriods(repo),
		application.NewListPeriodsByState(repo),
		application.NewGetPeriodByID(repo),
		application.NewUpdatePeriod(repo),
		application.NewClosePeriod(repo),
	)

	periods := router.Group("/periods")
	periods.PUT("/:id", handler.UpdatePeriod)

	body := delivery.UpdatePeriodRequest{
		Name: "2026-11",
	}

	bodyJSON, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/periods/1", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("PUT /periods/1 status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp delivery.PeriodResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "2026-11" {
		t.Errorf("PUT /periods/1 did not update correctly")
	}
}
