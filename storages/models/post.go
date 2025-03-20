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

type PostComment struct {
	ID              int64
	CommentId       string
	PostId          string
	ParentCommentId string
	ParentUserId    string
	Text            string
	IsDelete        int
	UserId          string
	CreatedTime     int64
	UpdatedTime     time.Time
	Status          int
	AppKey          string
}

type IPostCommentStorage interface {
	Create(item PostComment) error
	FindByIds(appkey string, commentIds []string) (map[string]*PostComment, error)
	QryPostComments(appkey, postId string, startTime, limit int64, isPosttive bool) ([]*PostComment, error)
}
