package model

// User represents an application user account.
type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateUsernameInput contains data for username update operation.
type UpdateUsernameInput struct {
	Username string `json:"username"`
}
