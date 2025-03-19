package dbs

import "time"

type PostCommentDao struct {
	ID              int64     `gorm:"primary_key"`
	CommentId       string    `gorm:"comment_id"`
	PostId          string    `gorm:"post_id"`
	ParentCommentId string    `gorm:"parent_comment_id"`
	ParentUserId    string    `gorm:"parent_user_id"`
	Content         string    `gorm:"content"`
	IsDelete        int       `gorm:"is_delete"`
	UserId          string    `gorm:"user_id"`
	CreatedTime     time.Time `gorm:"created_time"`
	UpdatedTime     time.Time `gorm:"updated_time"`
	Status          int       `gorm:"status"`
	AppKey          string    `gorm:"app_key"`
}

func (comment PostCommentDao) TableName() string {
	return "postcomments"
}
