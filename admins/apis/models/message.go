package models

type HisMsgs struct {
	Msgs []*HisMsg `json:"items"`
}

type HisMsg struct {
	Sender     *User  `json:"sender"`
	MsgId      string `json:"msg_id"`
	MsgTime    int64  `json:"msg_time"`
	MsgType    string `json:"msg_type"`
	MsgContent string `json:"msg_content"`
}
