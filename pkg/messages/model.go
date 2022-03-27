package messages

type MessageGet struct {
	Id          uint   `json:"messageId"`
	ThreadId    int    `json:"threadId"`
	FromId      int    `json:"fromId"`
	MessageText string `json:"messageText"`
	Date        int64  `json:"date"`
}

type MessagePost struct {
	Id          uint    `json:"messageId"`
	Thread      *Thread `json:"thread"`
	ThreadId    int     `json:"threadId"`
	FromId      int     `json:"fromId"`
	MessageText string  `json:"messageText"`
	Date        int64   `json:"date"`
}

type Thread struct {
	Id        int    `json:"threadId"`
	Name      string `json:"threadName"`
	CreatedAt int64
	Type      uint8 `json:"threadType"`
	CreatorId int   `json:"creatorId"`
	Users     []int `json:"participants"`
}