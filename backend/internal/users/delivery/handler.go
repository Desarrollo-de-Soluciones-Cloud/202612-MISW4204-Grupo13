package delivery

import (
	authDelivery "backend/internal/auth/delivery"
	authDomain "backend/internal/auth/domain"
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
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
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
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
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
		case isUserValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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
	currentUser, ok := authDelivery.GetCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, authDomain.ErrAuthTokenRequired)
		return
	}

	roleFilter := c.Query("role")
	if roleFilter != "" {
		if currentUser.GlobalRole == domain.RoleProfessor && !isProfessorAllowedRoleFilter(domain.UserRole(roleFilter)) {
			sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
			return
		}

		h.ListUsersByRole(c, roleFilter)
		return
	}

	if currentUser.GlobalRole == domain.RoleProfessor {
		sharedHelpers.RespondWithError(c, http.StatusForbidden, authDomain.ErrAuthForbidden)
		return
	}

	output, err := h.listUsers.Execute()
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	output, err := h.getUserByID.Execute(application.GetUserByIDInput{ID: id})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
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
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case errors.Is(err, domain.ErrUserEmailAlreadyInUse):
			sharedHelpers.RespondWithError(c, http.StatusConflict, err)
		case isUserValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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
	id, err := sharedHelpers.ParseResourceID(c.Param("id"))
	if err != nil {
		sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req ChangeUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.changeUserRole.Execute(application.ChangeUserRoleInput{
		ID:         id,
		GlobalRole: domain.UserRole(req.GlobalRole),
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			sharedHelpers.RespondWithError(c, http.StatusNotFound, err)
		case isUserValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
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

func isProfessorAllowedRoleFilter(role domain.UserRole) bool {
	return role == domain.RoleAssistant || role == domain.RoleMonitor
}
