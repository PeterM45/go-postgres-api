package database

type UserStore interface {
	CreateUser(username, email, password string) (*User, error)
	GetUserByID(id int) (*User, error)
	GetUsers() ([]User, error)
	VerifyUser(username, password string) (*User, error)
}

var _ UserStore = &DB{}
