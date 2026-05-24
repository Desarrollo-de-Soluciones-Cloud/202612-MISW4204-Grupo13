package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	nextCalled := false
	router := gin.New()
	router.Use(CORSMiddleware([]string{"https://frontend.example.com"}))
	router.GET("/", func(ctx *gin.Context) {
		nextCalled = true
		ctx.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://frontend.example.com")
	router.ServeHTTP(recorder, req)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "https://frontend.example.com" {
		t.Fatalf("expected allow origin header to be set")
	}
	if recorder.Header().Get("Vary") != "Origin" {
		t.Fatalf("expected vary origin header")
	}
	if !nextCalled {
		t.Fatalf("expected next handler to be called")
	}
}

func TestCORSMiddlewareAllowsAnyOriginWhenListIsEmpty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Origin", "https://random.example.com")

	CORSMiddleware(nil)(c)

	if recorder.Header().Get("Access-Control-Allow-Origin") != "https://random.example.com" {
		t.Fatalf("expected wildcard-by-empty-config behavior")
	}
}

func TestCORSMiddlewareRejectsUnknownOriginAndHandlesPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodOptions, "/", nil)
	c.Request.Header.Set("Origin", "https://blocked.example.com")

	CORSMiddleware([]string{"https://allowed.example.com"})(c)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("expected no allow origin header for blocked origin")
	}
}

func TestIsAllowedOrigin(t *testing.T) {
	allowed := map[string]struct{}{
		"https://frontend.example.com": {},
	}

	if !isAllowedOrigin("https://frontend.example.com", allowed) {
		t.Fatalf("expected origin to be allowed")
	}
	if isAllowedOrigin("https://blocked.example.com", allowed) {
		t.Fatalf("expected origin to be rejected")
	}
}
