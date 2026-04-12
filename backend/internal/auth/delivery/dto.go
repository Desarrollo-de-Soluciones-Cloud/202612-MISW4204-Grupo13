package delivery

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type AuthUserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GlobalRole string `json:"global_role"`
}

type SignInResponse struct {
	AccessToken string           `json:"access_token"`
	TokenType   string           `json:"token_type"`
	ExpiresIn   int64            `json:"expires_in"`
	User        AuthUserResponse `json:"user"`
}
