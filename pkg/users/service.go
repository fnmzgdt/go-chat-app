package users

import (
	"project/pkg/errors"
)

type Repository interface {
	ExecuteCreateUser(query string, values ...interface{})
	ExecuteGetUserByEmail(query string, values ...interface{}) (string, *errors.RestError)
}

type Service interface {
	CreateUser(user *User)
	GetUserByEmail(user *User) (string, *errors.RestError)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateUser(user *User) {
	query := `INSERT INTO users(first_name, last_name, email, password) VALUES(?, ?, ?, ?)`
	s.r.ExecuteCreateUser(query, user.FirstName, user.LastName, user.Email, user.Password)
}

func (s *service) GetUserByEmail(user *User) (string ,*errors.RestError) {
	query := `SELECT password FROM users WHERE email=?`
	password, error := s.r.ExecuteGetUserByEmail(query, user.Email)
	return password, error
}