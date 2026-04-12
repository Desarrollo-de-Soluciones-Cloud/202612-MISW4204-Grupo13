package delivery

import (
	sharedErrors "backend/internal/shared/errors"
	sharedHelpers "backend/internal/shared/helpers"
	"backend/internal/auth/application"
	"backend/internal/auth/domain"
	usersDomain "backend/internal/users/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const currentUserContextKey = "current_user"

type AuthHandler struct {
	signIn        *application.SignIn
	validateToken *application.ValidateToken
}

func NewAuthHandler(signIn *application.SignIn, validateToken *application.ValidateToken) *AuthHandler {
	return &AuthHandler{
		signIn:        signIn,
		validateToken: validateToken,
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sharedHelpers.RespondWithErrors(c, http.StatusBadRequest, mapBindingErrors(err))
		return
	}

	output, err := h.signIn.Execute(application.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			sharedHelpers.RespondWithError(c, http.StatusUnauthorized, err)
		case isAuthValidationError(err):
			sharedHelpers.RespondWithError(c, http.StatusBadRequest, err)
		default:
			sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, SignInResponse{
		AccessToken: output.AccessToken,
		TokenType:   string(output.TokenType),
		ExpiresIn:   output.ExpiresIn,
		User: AuthUserResponse{
			ID:         output.User.ID,
			Name:       output.User.Name,
			Email:      output.User.Email,
			GlobalRole: string(output.User.GlobalRole),
		},
	})
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	currentUser, ok := getCurrentUser(c)
	if !ok {
		sharedHelpers.RespondWithError(c, http.StatusUnauthorized, domain.ErrAuthTokenRequired)
		return
	}

	c.JSON(http.StatusOK, AuthUserResponse{
		ID:         currentUser.ID,
		Name:       currentUser.Name,
		Email:      currentUser.Email,
		GlobalRole: string(currentUser.GlobalRole),
	})
}

func (h *AuthHandler) RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			sharedHelpers.RespondWithError(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		output, err := h.validateToken.Execute(application.ValidateTokenInput{Token: token})
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrAuthTokenRequired),
				errors.Is(err, domain.ErrAuthTokenInvalid),
				errors.Is(err, domain.ErrAuthTokenExpired):
				sharedHelpers.RespondWithError(c, http.StatusUnauthorized, err)
			default:
				sharedHelpers.RespondWithError(c, http.StatusInternalServerError, sharedErrors.ErrInternalServerError)
			}
			c.Abort()
			return
		}

		c.Set(currentUserContextKey, domain.AuthenticatedUser{
			ID:         output.ID,
			Name:       output.Name,
			Email:      output.Email,
			GlobalRole: usersDomain.UserRole(output.GlobalRole),
		})
		c.Next()
	}
}

func (h *AuthHandler) RequireRoles(roles ...usersDomain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, ok := getCurrentUser(c)
		if !ok {
			sharedHelpers.RespondWithError(c, http.StatusUnauthorized, domain.ErrAuthTokenRequired)
			c.Abort()
			return
		}

		for _, role := range roles {
			if currentUser.GlobalRole == role {
				c.Next()
				return
			}
		}

		sharedHelpers.RespondWithError(c, http.StatusForbidden, domain.ErrAuthForbidden)
		c.Abort()
	}
}

func getCurrentUser(c *gin.Context) (domain.AuthenticatedUser, bool) {
	rawCurrentUser, exists := c.Get(currentUserContextKey)
	if !exists {
		return domain.AuthenticatedUser{}, false
	}

	currentUser, ok := rawCurrentUser.(domain.AuthenticatedUser)
	if !ok {
		return domain.AuthenticatedUser{}, false
	}

	return currentUser, true
}

func isAuthValidationError(err error) bool {
	return errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrAuthEmailRequired) ||
		errors.Is(err, domain.ErrAuthEmailInvalid) ||
		errors.Is(err, domain.ErrAuthPasswordRequired) ||
		errors.Is(err, domain.ErrAuthPasswordTooShort) ||
		errors.Is(err, domain.ErrAuthPasswordTooLong)
}
