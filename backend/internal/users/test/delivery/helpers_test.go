package delivery_test

import (
	authDomain "backend/internal/auth/domain"
	applicationpkg "backend/internal/users/application"
	deliverypkg "backend/internal/users/delivery"
	usersDomain "backend/internal/users/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestListUsersByRoleBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := newMockUserRepository()
	handler := deliverypkg.NewUserHandler(
		applicationpkg.NewCreateUser(repo),
		applicationpkg.NewListUsers(repo),
		applicationpkg.NewListUsersByRole(repo),
		applicationpkg.NewGetUserByID(repo),
		applicationpkg.NewUpdateUser(repo),
		applicationpkg.NewChangeUserRole(repo),
	)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users?role=invalid", nil)
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.ListUsersByRole(c, "invalid")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
