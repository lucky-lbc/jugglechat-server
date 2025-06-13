package models

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
