package model

type User struct {
	Password string
}
type CreateUserData struct {
	Username string
}
type ChangeUser struct {
	NewUsername string
}
