package delivery

import (
	"backend/internal/users/application"
	"backend/internal/users/domain"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    createUser *application.CreateUser
    listUsers  *application.ListUsers
}

func NewUserHandler(createUser *application.CreateUser, listUsers *application.ListUsers) *UserHandler {
    return &UserHandler{
        createUser: createUser,
        listUsers:  listUsers,
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    input := application.CreateUserInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }
    output, err := h.createUser.Execute(input)
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrUserAlreadyExists):
            c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
        case errors.Is(err, domain.ErrInvalidInput):
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        }
        return
    }
    c.JSON(http.StatusCreated, CreateUserResponse{
        ID:    output.ID,
        Name:  output.Name,
        Email: output.Email,
    })
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    output, err := h.listUsers.Execute()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }
    users := make([]UserResponse, len(output.Users))
    for i, u := range output.Users {
        users[i] = UserResponse{
            ID:    u.ID,
            Name:  u.Name,
            Email: u.Email,
        }
    }
    c.JSON(http.StatusOK, ListUsersResponse{Users: users})
}