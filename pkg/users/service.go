package users

type Repository interface {
	ExecuteQuery(query string, values ...interface{}) 
}

type Service interface {
	AddUser(user *User)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddUser(user *User) {
	query := `INSERT INTO users(first_name, last_name, email) VALUES(?, ?, ?)`
	s.r.ExecuteQuery(query, user.FirstName, user.LastName, user.Email)
}