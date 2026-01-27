package response

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   int8   `json:"status"`
}

type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt int64       `json:"expires_at"`
	User      *UserResponse `json:"user"`
}

type ListUsersResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}
