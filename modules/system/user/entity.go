package user

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	RoleID   uint   `json:"role_id"`
}

type UpdateUserRequest struct {
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	RoleID   uint   `json:"role_id,omitempty"`
	Status   *int   `json:"status,omitempty"`
}
