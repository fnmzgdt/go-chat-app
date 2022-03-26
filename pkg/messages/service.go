package messages

import (
	"database/sql"
)

type Service interface {
	GetUserThread(user1id int, user2id int) (int64, error)
	CreateThread(thread *Thread,  date int64) (sql.Result, error)
	CreateUserThread(userid int, threadid int64, date int64, seen uint) (sql.Result, error)
	CreateMessage(message *Message) (sql.Result, error)
}

type Repository interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	ExecuteGetUserThread(query string, values ...interface{}) (int64, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetUserThread(user1id int, user2id int) (int64, error) {
	query := "SELECT ut1.thread_id FROM users_threads ut1 JOIN users_threads ut2 ON ut1.thread_id = ut2.thread_id JOIN threads t ON ut1.thread_id = t.id WHERE  ut1.user_id = ? AND ut2.user_id = ? AND t.type = 0;"
	result, err := s.r.ExecuteGetUserThread(query, user1id, user2id)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *service) CreateThread(thread *Thread, date int64) (sql.Result, error) {
	query := "INSERT INTO threads(name, created_at, type) VALUES(?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, thread.Name, date, thread.Type)
	if err != nil {
		return nil, err
	}
	return result, nil
} 

func (s *service) CreateUserThread(userid int, threadid int64, date int64, seen uint) (sql.Result, error) {
	query := "INSERT INTO users_threads(user_id, thread_id, added_by, date, seen) VALUES (?, ?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, userid, threadid, sql.NullInt64{}, date, seen)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) CreateMessage(message *Message) (sql.Result, error) {
	query := "INSERT INTO messages(thread_id, user_id, date, text) VALUES (?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, message.ThreadId, message.FromId, message.Date, message.MessageText)
	if err != nil {
		return nil, err
	}
	return result, nil
}