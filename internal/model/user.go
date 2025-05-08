package model

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

type UpdateUsernameInput struct {
	Username string `json:"username"`
}
