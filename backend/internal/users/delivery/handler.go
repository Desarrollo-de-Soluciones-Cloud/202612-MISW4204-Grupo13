package delivery

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	createUser      *application.CreateUser
	listUsers       *application.ListUsers
	listUsersByRole *application.ListUsersByRole
	getUserByID     *application.GetUserByID
	updateUser      *application.UpdateUser
	changeUserRole  *application.ChangeUserRole
}

func NewUserHandler(
	createUser *application.CreateUser,
	listUsers *application.ListUsers,
	listUsersByRole *application.ListUsersByRole,
	getUserByID *application.GetUserByID,
	updateUser *application.UpdateUser,
	changeUserRole *application.ChangeUserRole,
) *UserHandler {
	return &UserHandler{
		createUser:      createUser,
		listUsers:       listUsers,
		listUsersByRole: listUsersByRole,
		getUserByID:     getUserByID,
		updateUser:      updateUser,
		changeUserRole:  changeUserRole,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := application.CreateUserInput{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		GlobalRole: domain.UserRole(req.GlobalRole),
	}
	output, err := h.createUser.Execute(input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserEmailAlreadyInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case isUserValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, CreateUserResponse{
		ID:         output.ID,
		Name:       output.Name,
		Email:      output.Email,
		GlobalRole: string(output.GlobalRole),
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	roleFilter := c.Query("role")
	if roleFilter != "" {
		h.ListUsersByRole(c, roleFilter)
		return
	}

	output, err := h.listUsers.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	users := make([]UserResponse, len(output.Users))
	for i, u := range output.Users {
		users[i] = UserResponse{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			GlobalRole: string(u.GlobalRole),
		}
	}
	c.JSON(http.StatusOK, ListUsersResponse{Users: users})
}

func (h *UserHandler) ListUsersByRole(c *gin.Context, rawRole string) {
	output, err := h.listUsersByRole.Execute(application.ListUsersByRoleInput{
		GlobalRole: domain.UserRole(rawRole),
	})
	if err != nil {
		switch {
		case isUserValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	users := make([]UserResponse, len(output.Users))
	for i, u := range output.Users {
		users[i] = UserResponse{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			GlobalRole: string(u.GlobalRole),
		}
	}

	c.JSON(http.StatusOK, ListUsersResponse{Users: users})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.getUserByID.Execute(application.GetUserByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:         output.ID,
		Name:       output.Name,
		Email:      output.Email,
		GlobalRole: string(output.GlobalRole),
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.updateUser.Execute(application.UpdateUserInput{
		ID:         id,
		Name:       req.Name,
		Email:      req.Email,
		GlobalRole: domain.UserRole(req.GlobalRole),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, domain.ErrUserEmailAlreadyInUse):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case isUserValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:         output.ID,
		Name:       output.Name,
		Email:      output.Email,
		GlobalRole: string(output.GlobalRole),
	})
}

func (h *UserHandler) ChangeUserRole(c *gin.Context) {
	id, err := parseUserIDParam(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req ChangeUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.changeUserRole.Execute(application.ChangeUserRoleInput{
		ID:         id,
		GlobalRole: domain.UserRole(req.GlobalRole),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case isUserValidationError(err):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:         output.ID,
		Name:       output.Name,
		Email:      output.Email,
		GlobalRole: string(output.GlobalRole),
	})
}

func isUserValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrUserNameRequired) ||
		errors.Is(err, domain.ErrUserNameTooShort) ||
		errors.Is(err, domain.ErrUserNameTooLong) ||
		errors.Is(err, domain.ErrUserEmailRequired) ||
		errors.Is(err, domain.ErrUserEmailInvalid) ||
		errors.Is(err, domain.ErrUserPasswordRequired) ||
		errors.Is(err, domain.ErrUserPasswordTooShort) ||
		errors.Is(err, domain.ErrUserPasswordTooLong) ||
		errors.Is(err, domain.ErrUserRoleRequired) ||
		errors.Is(err, domain.ErrUserRoleInvalid)
}
