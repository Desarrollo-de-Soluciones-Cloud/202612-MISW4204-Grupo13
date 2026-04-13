package delivery_test

import (
	"backend/internal/weeks/application"
	"backend/internal/weeks/delivery"
	"backend/internal/weeks/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockWeekRepository struct {
	weeks []domain.Week
}

func (m *MockWeekRepository) CreateMany(weeks []domain.Week) error {
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
	return false, nil
}

func setupTestHandler(repo domain.WeekRepository) *delivery.WeekHandler {
	return delivery.NewWeekHandler(
		application.NewListWeeksByPeriod(repo),
		application.NewGetWeekByPeriodAndNumber(repo),
	)
}

func TestListWeeksByPeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
			{ID: 2, PeriodID: 1, Number: 2, InitialDate: "2026-01-19", FinalDate: "2026-01-25"},
		},
	}
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/periods/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "1"})

	handler.ListWeeksByPeriod(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp delivery.ListWeeksResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if len(resp.Weeks) != 2 {
		t.Fatalf("expected 2 weeks, got %d", len(resp.Weeks))
	}
}

func TestGetWeekByPeriodAndNumberSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockWeekRepository{
		weeks: []domain.Week{
			{ID: 1, PeriodID: 1, Number: 1, InitialDate: "2026-01-12", FinalDate: "2026-01-18"},
		},
	}
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/1/periods/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "1"})
	c.Params = append(c.Params, gin.Param{Key: "number", Value: "1"})

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp delivery.WeekResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if resp.Number != 1 {
		t.Fatalf("expected week number 1, got %d", resp.Number)
	}
}

func TestGetWeekByPeriodAndNumberNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &MockWeekRepository{}
	handler := setupTestHandler(repo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/1/periods/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "1"})
	c.Params = append(c.Params, gin.Param{Key: "number", Value: "1"})

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestListWeeksByPeriodInvalidPeriodID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler(&MockWeekRepository{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/periods/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "invalid"})

	handler.ListWeeksByPeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetWeekByPeriodAndNumberInvalidPeriodID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler(&MockWeekRepository{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/1/periods/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "number", Value: "1"})
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "invalid"})

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetWeekByPeriodAndNumberInvalidNumber(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler(&MockWeekRepository{})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/weeks/invalid/periods/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{Key: "number", Value: "invalid"})
	c.Params = append(c.Params, gin.Param{Key: "periodId", Value: "1"})

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
