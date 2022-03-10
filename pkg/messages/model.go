package messages

type Message struct {
	ThreadId string `json:"threadId"`
	FromId string `json:"fromId"`
	ToId string `json:"toId"`
	FromUsername string `json:"fromUsername"`
	ToUsername string `json:"toUsername"`
	Message string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

type Thread struct {

}