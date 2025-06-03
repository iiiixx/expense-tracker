package model

// LoginInput contains user credentials for authentication.
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
