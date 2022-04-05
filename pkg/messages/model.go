package messages

import (
	"errors"
	"strings"
)

type MessageGet struct {
	Id          uint   `json:"messageId"`
	ThreadId    int    `json:"threadId"`
	FromId      int    `json:"fromId"`
	MessageText string `json:"messageText"`
	Date        int64  `json:"date"`
}

func (m MessagePost) validate() error {
	if len(strings.Trim(m.MessageText, " ")) == 0 {
		return errors.New("you cannot send empty messages")
	}
	return nil
}

type MessagePost struct {
	Id          uint    `json:"messageId,omitempty"`
	Thread      *Thread `json:"thread,omitempty"`
	ThreadId    int     `json:"threadId,omitempty"`
	FromId      int     `json:"fromId,omitempty"`
	MessageText string  `json:"messageText,omitempty"`
	Date        int64   `json:"date,omitempty"`
}

type Thread struct {
	Id        int    `json:"threadId,omitempty"`
	Name      string `json:"threadName,omitempty"`
	CreatedAt int64
	Type      uint8 `json:"threadType,omitempty"`
	CreatorId int   `json:"creatorId,omitempty"`
	Users     []int `json:"participants,omitempty"`
}

type ThreadGet struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Seen 		uint8 `json:"seen"`
	LastSender string `json:"lastSender"`
	LastMessage string `json:"lastMessage"`
	LastUpdated int64 `json:"lastUpdated"`
	Type      uint8 `json:"type"`
	Users     string `json:"participants"`
}