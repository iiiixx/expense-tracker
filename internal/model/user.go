package model

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUsernameInput struct {
	Username string `json:"username"`
}
