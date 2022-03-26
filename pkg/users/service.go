package users

import "database/sql"


type Repository interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	ExecuteGetUserByEmail(query string, values ...interface{}) (string, error)
}

type Service interface {
	CreateUser(user *User) (sql.Result, error)
	GetUserByEmail(user *User) (string, error)
}

type service struct {
	mysql Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateUser(user *User) (sql.Result, error) {
	query := `INSERT INTO users (username, email, password) VALUES(?, ?, ?);`
	result, err := s.mysql.ExecuteQuery(query, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetUserByEmail(user *User) (string, error) {
	query := `SELECT CAST(password AS CHAR) FROM users WHERE email=?;`
	password, error := s.mysql.ExecuteGetUserByEmail(query, user.Email)
	if error != nil {
		return "", error
	}
	return password, nil
} 