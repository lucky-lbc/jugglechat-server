package models

type Posts struct {
	Items      []*Post `json:"items"`
	IsFinished bool    `json:"is_finished"`
}

type Post struct {
	PostId      string                 `json:"post_id"`
	Content     *PostContent           `json:"content"`
	UserInfo    *UserObj               `json:"user_info"`
	CreatedTime int64                  `json:"created_time"`
	UpdatedTime int64                  `json:"updated_time"`
	Reactions   map[string][]*Reaction `json:"reactions"`
	TopComments []*PostComment         `json:"top_comments"`
}

type Reaction struct {
	Value    string   `json:"value"`
	UserInfo *UserObj `json:"user_info"`
}

type PostContent struct {
	Text   string              `json:"text"`
	Images []*PostContentImage `json:"images"`
	Video  *PostContentVideo   `json:"video"`
}

type PostContentImage struct {
	Url string `json:"url"`
}

type PostContentVideo struct {
	Url string `json:"url"`
}

type PostComment struct {
	CommentId       string   `json:"comment_id"`
	PostId          string   `json:"post_id"`
	ParentCommentId string   `json:"parent_comment_id"`
	Text            string   `json:"text"`
	ParentUserId    string   `json:"parent_user_id,omitempty"`
	ParentUserInfo  *UserObj `json:"parent_user_info"`
	UserInfo        *UserObj `json:"user_info"`
	CreatedTime     int64    `json:"created_time"`
	UpdatedTime     int64    `json:"updated_time"`
}

type PostComments struct {
	Items      []*PostComment `json:"items"`
	IsFinished bool           `json:"is_finished"`
}
