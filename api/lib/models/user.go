package models

// User a model for defining an user
type User struct {
	ID       int64
	Email    string
	PassHash []byte
}
