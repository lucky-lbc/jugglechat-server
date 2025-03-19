package models

import "time"

type Post struct {
	ID           int64
	PostId       string
	Title        string
	Content      []byte
	ContentBrief string
	PostExset    []byte
	IsDelete     int
	UserId       string
	CreatedTime  int64
	UpdatedTime  time.Time
	Status       int
	AppKey       string
}

type IPostStorage interface {
	Create(item Post) error
	FindById(appkey, postId string) (*Post, error)
	FindByIds(appkey string, postIds []string) (map[string]*Post, error)
	QryPosts(appkey string, startTime, limit int64, isPositive bool) ([]*Post, error)
}
