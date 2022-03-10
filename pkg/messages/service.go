package messages

type Service interface {
	SendMessage(message *Message)
	CreateThread(thread *Thread)
}

type Repository interface {
	ExecuteCreateMessage(query string, threadId string, createdAt string, message string, fromId string, toId string, fromUsername string, toUsername string)
	ExecuteCreateThread(query string, values ...interface{})
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) SendMessage(message *Message) {
	query := `INSERT INTO message_by_thread_id(thread_id, message_id, created_time, is_viewed, message, user_data) VALUES(?, now() , ?, false, ?, {from_id: ?, to_id: ?, from_username: ?, to_username: ?});`
	s.r.ExecuteCreateMessage(query, message.ThreadId, message.CreatedAt, message.Message, message.FromId, message.ToId, message.FromUsername, message.ToUsername)
}

func (s *service) CreateThread(thread *Thread) {

}