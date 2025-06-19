package models

import "time"

type Feedback struct {
	AppKey      string
	UserId      string
	Category    string
	Content     []byte
	UpdatedTime time.Time
	CreatedTime time.Time
}

type IFeedbackStorage interface {
	Create(item Feedback) error
}
