package messages

type Message struct {
	ThreadId    int64    `json:"threadId"`
	FromId      int    `json:"fromId"`
	ToId        int    `json:"toId"`
	MessageText string `json:"messageText"`
	Date        int64
}

type Thread struct {
	Name      string `json:"threadName"`
	Type      uint8 `json:"threadType"`
}