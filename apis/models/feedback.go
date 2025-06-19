package models

type Feedback struct {
	Category string   `json:"category"`
	Text     string   `json:"text"`
	Images   []string `json:"images"`
	Videos   []string `json:"videos"`

	AppKey      string   `json:"app_key"`
	User        *UserObj `json:"user"`
	CreatedTime int64    `json:"created_time"`
	UpdatedTime int64    `json:"updated_time"`
}
