package apimodels

type BotType int

var (
	BotType_Default     BotType = 0
	BotType_Custom      BotType = 1
	BotType_Dify        BotType = 2
	BotType_Coze        BotType = 3
	BotType_Minmax      BotType = 4
	BotType_SiliconFlow BotType = 5
)

type BotMsg struct {
	SenderId    string        `json:"sender_id"`
	BotId       string        `json:"bot_id"`
	ChannelType int           `json:"channel_type"`
	Stream      bool          `json:"stream"`
	Messages    []*BotMsgItem `json:"messages"`
}

type BotMsgItem struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type BotResponsePartData struct {
	Id             string `json:"id"`
	ConversationId string `json:"conversation_id"`
	Type           string `json:"type"`
	BotId          string `json:"bot_id"`
	Content        string `json:"content"`
	ContentType    string `json:"content_type"`
	SectionId      string `json:"section_id"`
	CreatedTime    int64  `json:"created_time"`
}

type AiBotInfos struct {
	Items  []*AiBotInfo `json:"items"`
	Offset string       `json:"offset"`
}

type AiBotInfo struct {
	BotId    string `json:"bot_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	BotType  int32  `json:"bot_type"`
}

type TelegramBot struct {
	BotId       string `json:"bot_id"`
	BotName     string `json:"bot_name"`
	BotToken    string `json:"bot_token"`
	CreatedTime int64  `json:"created_time"`
}

type TelegramBots struct {
	Items  []*TelegramBot `json:"items"`
	Offset string         `json:"offset"`
}

type TelegramBotIds struct {
	BotIds []string `json:"bot_ids"`
}
