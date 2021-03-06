package messages

import (
	"database/sql"
	"fmt"
)

type Service interface {
	GetUserThread(user1id int, user2id int) (int64, error)
	CreateThread(thread *ThreadPost) (sql.Result, error)
	CreateUserThread(userid int, threadid int64, date int64, seen uint8) (sql.Result, error)
	CreateMessage(message *MessagePost) (sql.Result, error)
	GetUsersInThread(threadId int) ([]int, error)
	UpdateUserThread(userid int, threadId int, userIds []int) (sql.Result, error)
	GetMessagesFromThread(threadid int) ([]MessageGet, error)
	getLatestThreads(userid int) ([]ThreadGet, error)
}

type Repository interface {
	ExecuteQuery(query string, values ...interface{}) (sql.Result, error)
	ExecuteUpdateUserThreadQuery(query string, values []int) (sql.Result, error)
	ExecuteGetUserThread(query string, values ...interface{}) (int64, error)
	ExecuteGetUsersInThread(query string, threadId int) ([]int, error)
	ExecuteGetMessagesFromThread(query string, threadid int) ([]MessageGet, error)
	ExecuteGetLatestThreads(query string, userid int) ([]ThreadGet, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

//not used anywhere rn
func (s *service) GetUserThread(user1id int, user2id int) (int64, error) {
	query := "SELECT ut1.thread_id FROM users_threads ut1 JOIN users_threads ut2 ON ut1.thread_id = ut2.thread_id JOIN threads t ON ut1.thread_id = t.id WHERE  ut1.user_id = ? AND ut2.user_id = ? AND t.type = 0;"
	result, err := s.r.ExecuteGetUserThread(query, user1id, user2id)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (s *service) CreateThread(thread *ThreadPost) (sql.Result, error) {
	query := "INSERT INTO threads(name, created_at, type, created_by) VALUES(NULLIF(?, ''), from_unixtime(?), ?, NULLIF(?, 0));"

	result, err := s.r.ExecuteQuery(query, thread.Name, thread.CreatedAt, thread.Type, thread.CreatorId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//not used anywhere rn
func (s *service) CreateUserThread(userid int, threadid int64, date int64, seen uint8) (sql.Result, error) {
	query := "INSERT INTO users_threads(user_id, thread_id, added_by, date, seen) VALUES (?, ?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, userid, threadid, sql.NullInt64{}, date, seen)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetUsersInThread(threadId int) ([]int, error) {
	query := "SELECT user_id FROM users_threads WHERE thread_id = ?;"
	result, err := s.r.ExecuteGetUsersInThread(query, threadId)
	if err != nil {
		return []int{}, err
	}

	return result, nil
}

func (s *service) UpdateUserThread(userid int, threadId int, userIds []int) (sql.Result, error) {
	queryTemplate := "UPDATE users_threads ut JOIN (%s) temp ON temp.id = ut.user_id AND temp.thread_id = ut.thread_id SET ut.seen = temp.seen;"
	var queryParams []int
	innerQuery := ""
	for i, userId := range userIds {

		if i == 0 {
			innerQuery += "SELECT ? as id, ? as thread_id, ? as seen "
		} else {
			innerQuery += "UNION ALL SELECT ?, ?, ? "
		}

		seen := 0
		if userIds[i] == userid {
			seen = 1
		}

		queryParams = append(queryParams, userId, threadId, seen)
	}

	query := fmt.Sprintf(queryTemplate, innerQuery)

	result, err := s.r.ExecuteUpdateUserThreadQuery(query, queryParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) CreateMessage(message *MessagePost) (sql.Result, error) {
	query := "INSERT INTO messages(thread_id, user_id, date, text) VALUES (?, ?, from_unixtime(?), ?);"
	result, err := s.r.ExecuteQuery(query, message.ThreadId, message.FromId, message.Date, message.MessageText)
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

func (s *service) getLatestThreads(userid int) ([]ThreadGet, error) {
	query := "SELECT b.`id` AS `id`, b.`name` AS `name`, a.`seen`, d.`username` AS `lastSender`, c.`text` AS `lastMessage`, UNIX_TIMESTAMP(c.`date`) AS `lastUpdated`, b.`type` AS `type`, f.`participants` FROM `users_threads` a JOIN `threads` b ON a.`thread_id` = b.`id` JOIN (SELECT messages.`id`, `thread_id`, `text`, `date`, `user_id` AS `sender_id` FROM messages JOIN (SELECT Max(id) AS `id` FROM messages WHERE thread_id IN (SELECT thread_id FROM users_threads WHERE user_id = ?) GROUP BY thread_id ORDER BY id DESC LIMIT 20) b ON messages.id = b.id) c ON b.`id` = c.`thread_id` JOIN users d ON d.id = c.sender_id JOIN (SELECT thread_id, Group_concat(b.username, '') AS `participants` FROM users_threads a JOIN users b ON a.user_id = b.id WHERE thread_id IN (SELECT thread_id FROM users_threads WHERE user_id = ?) AND user_id != ? GROUP BY thread_id) f ON f.thread_id = a.thread_id WHERE a.`user_id` = ? ORDER BY c.`id` DESC;"
	result, err := s.r.ExecuteGetLatestThreads(query, userid)
	if err != nil {
		return nil, err
	}
	return result, nil
}
