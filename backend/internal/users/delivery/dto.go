package delivery

type CreateUserRequest struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8,max=72"`
	GlobalRole string `json:"global_role" binding:"required"`
}

type CreateUserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name       string `json:"name" binding:"required,min=3,max=100"`
	Email      string `json:"email" binding:"required,email"`
	GlobalRole string `json:"global_role" binding:"required"`
}

type ChangeUserRoleRequest struct {
	GlobalRole string `json:"global_role" binding:"required"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
