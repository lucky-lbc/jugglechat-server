package apimodels

type Posts struct {
	Items      []*Post `json:"items"`
	IsFinished bool    `json:"is_finished"`
}

type Post struct {
	PostId      string               `json:"post_id"`
	Content     *PostContent         `json:"content"`
	UserInfo    *UserObj             `json:"user_info"`
	CreatedTime int64                `json:"created_time"`
	UpdatedTime int64                `json:"updated_time"`
	Reactions   map[string]*Reaction `json:"reactions"`
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
	PostId          string   `json:"post_id"`
	CommentId       string   `json:"comment_id"`
	ParentCommentId string   `json:"parent_comment_id"`
	Content         string   `json:"content"`
	ParentUserInfo  *UserObj `json:"parent_user_info"`
	UserInfo        *UserObj `json:"user_info"`
	CreatedTime     int64    `json:"created_time"`
	UpdatedTime     int64    `json:"updated_time"`
}
