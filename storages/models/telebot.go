package models

import "time"

type TeleBot struct {
	ID          int64
	UserId      string
	BotName     string
	BotToken    string
	Status      int
	CreatedTime time.Time
	AppKey      string
}

type ITeleBotStorage interface {
	Create(item TeleBot) (int64, error)
	FindById(id int64, appkey, userId string) (*TeleBot, error)
	BatchDel(appkey, userId string, botIds []int64) error
	QryTeleBots(appkey, userId string, startId, limit int64) ([]*TeleBot, error)
}
