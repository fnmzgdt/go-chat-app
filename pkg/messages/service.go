package messages

import (
	"database/sql"
)

type Service interface {
	GetUserThread(user1id int, user2id int) (int64, error)
	CreateThread(thread *Thread) (sql.Result, error)
	CreateUserThread(userid int, threadid int64, date int64, seen uint8) (sql.Result, error)
	CreateMessage(message *MessagePost) (sql.Result, error)
	UpdateUserThread(userid int, threadid int, seen uint8) (sql.Result, error)
	GetMessagesFromThread(threadid int) ([]MessageGet, error)
}

type Repository interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	ExecuteGetUserThread(query string, values ...interface{}) (int64, error)
	ExecuteGetMessagesFromThread(query string, threadid int) ([]MessageGet, error)
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

func (s *service) CreateThread(thread *Thread) (sql.Result, error) {
	query := "INSERT INTO threads(name, created_at, type, created_by) VALUES(NULLIF(?, ''), from_unixtime(?), ?, NULLIF(?, 0));"
	
	result, err := s.r.ExecuteQuery(query, thread.Name, thread.CreatedAt, thread.Type, thread.CreatorId)
	if err != nil {
		return nil, err
	}
	return result, nil
} 

func (s *service) CreateUserThread(userid int, threadid int64, date int64, seen uint8) (sql.Result, error) {
	query := "INSERT INTO users_threads(user_id, thread_id, added_by, date, seen) VALUES (?, ?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, userid, threadid, sql.NullInt64{}, date, seen)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) UpdateUserThread(userid int, threadid int, seen uint8) (sql.Result, error) {
	query := "UPDATE users_threads SET seen = ? WHERE user_id = ? AND thread_id = ?"
	result, err := s.r.ExecuteQuery(query, seen, userid, threadid)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) CreateMessage(message *MessagePost) (sql.Result, error) {
	query := "INSERT INTO messages(thread_id, user_id, date, text) VALUES (?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, message.Thread.Id, message.FromId, message.Date, message.MessageText)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetMessagesFromThread(threadid int) ([]MessageGet, error) {
	query := "SELECT messages.id as messageId, thread_id as threadId, user_id as fromId, UNIX_TIMESTAMP(date) as date, text as messageText FROM messages JOIN users u ON messages.user_id = u.id WHERE thread_id = ? ORDER BY messages.id DESC;"
	result, err := s.r.ExecuteGetMessagesFromThread(query, threadid)
	if err != nil {
		return nil, err
	}
	return result, nil
}