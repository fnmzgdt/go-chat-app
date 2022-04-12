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

type MessagePost struct {
	ToUserIds   []int  `json:"thread,omitempty"`
	ThreadId    int    `json:"threadId,omitempty"`
	FromId      int    `json:"fromId,omitempty"`
	MessageText string `json:"messageText,omitempty"`
	Date        int    `json:"date,omitempty"`
}

func (m MessagePost) validate() error {
	if len(strings.Trim(m.MessageText, " ")) == 0 {
		return errors.New("Cannot send an empty message.")
	}

	if len(m.ToUserIds) == 0 {
		return errors.New("Cannot send message to an empty thread.")
	}

	if m.Date == 0 {
		return errors.New("Date field is required to have a value")
	}

	return nil
}

type ThreadPost struct {
	Name      string `json:"threadName,omitempty"`
	CreatedAt int
	Type      uint8 `json:"threadType,omitempty"`
	CreatorId int   `json:"creatorId,omitempty"`
	Users     []int `json:"participants,omitempty"`
}

func (t *ThreadPost) checkFields() error {
	if t.CreatedAt == 0 {
		return errors.New("CreatedAt field is required to have a value")
	}

	if t.CreatorId == 0 {
		return errors.New("CreatorId field is required to have a value")
	}

	if len(t.Users) == 0 {
		return errors.New("Users field is required to have a value")
	}

	return nil
}

type ThreadGet struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Seen        uint8  `json:"seen"`
	LastSender  string `json:"lastSender"`
	LastMessage string `json:"lastMessage"`
	LastUpdated int64  `json:"lastUpdated"`
	Type        uint8  `json:"type"`
	Users       string `json:"participants"`
}
