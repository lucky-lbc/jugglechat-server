package models

type RecallMsgReq struct {
	FromId      string            `json:"from_id"`
	TargetId    string            `json:"target_id"`
	ChannelType int32             `json:"channel_type"`
	SubChannel  string            `json:"sub_channel"`
	MsgId       string            `json:"msg_id"`
	MsgTime     int64             `json:"msg_time"`
	Exts        map[string]string `json:"exts"`
}

type DelHisMsgsReq struct {
	FromId      string       `json:"from_id"`
	TargetId    string       `json:"target_id"`
	ChannelType int32        `json:"channel_type"`
	SubChannel  string       `json:"sub_channel"`
	Msgs        []*SimpleMsg `json:"msgs"`
}

type SimpleMsg struct {
	MsgId        string `json:"msg_id"`
	MsgTime      int64  `json:"msg_time"`
	MsgReadIndex int64  `json:"msg_read_index"`
}
