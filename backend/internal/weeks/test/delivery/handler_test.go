package delivery_test

import (
	weeksApplication "backend/internal/weeks/application"
	weeksDelivery "backend/internal/weeks/delivery"
	weeksDomain "backend/internal/weeks/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type weekRepoStub struct {
	weeks []weeksDomain.Week
}

func (r *weekRepoStub) CreateMany(weeks []weeksDomain.Week) error { return nil }

func (r *weekRepoStub) FindAllByPeriodID(periodID uint) ([]weeksDomain.Week, error) {
	result := make([]weeksDomain.Week, 0)
	for _, week := range r.weeks {
		if week.PeriodID == periodID {
			result = append(result, week)
		}
	}
	return result, nil
}

func (r *weekRepoStub) FindByPeriodIDAndNumber(periodID uint, number int) (*weeksDomain.Week, error) {
	for _, week := range r.weeks {
		if week.PeriodID == periodID && week.Number == number {
			copy := week
			return &copy, nil
		}
	}
	return nil, weeksDomain.ErrWeekNotFound
}

func (r *weekRepoStub) FindByPeriodIDAndStartDate(periodID uint, startDate string) (*weeksDomain.Week, error) {
	return nil, weeksDomain.ErrWeekNotFound
}

func (r *weekRepoStub) ExistsByPeriodID(periodID uint) (bool, error) { return false, nil }

func newWeekHandlerForTest() *weeksDelivery.WeekHandler {
	repo := &weekRepoStub{
		weeks: []weeksDomain.Week{
			{ID: 1, PeriodID: 9, Number: 1, InitialDate: "2026-01-20", FinalDate: "2026-01-26"},
			{ID: 2, PeriodID: 9, Number: 2, InitialDate: "2026-01-27", FinalDate: "2026-02-02"},
		},
	}

	return weeksDelivery.NewWeekHandler(
		weeksApplication.NewListWeeksByPeriod(repo),
		weeksApplication.NewGetWeekByPeriodAndNumber(repo),
	)
}

func TestListWeeksByPeriodSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWeekHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weeks/periods/9", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "periodId", Value: "9"}}

	handler.ListWeeksByPeriod(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp weeksDelivery.ListWeeksResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Weeks) != 2 {
		t.Fatalf("expected 2 weeks, got %d", len(resp.Weeks))
	}
}

func TestListWeeksByPeriodBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWeekHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weeks/periods/bad", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "periodId", Value: "bad"}}

	handler.ListWeeksByPeriod(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetWeekByPeriodAndNumberSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWeekHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weeks/periods/9/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "periodId", Value: "9"}, {Key: "number", Value: "1"}}

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetWeekByPeriodAndNumberNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newWeekHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/weeks/periods/9/99", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "periodId", Value: "9"}, {Key: "number", Value: "99"}}

	handler.GetWeekByPeriodAndNumber(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
