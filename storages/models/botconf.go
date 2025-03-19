package models

type BotStatus int

var (
	BotStatus_Disable BotStatus = 0
	BotStatus_Enable  BotStatus = 1
)

type BotConf struct {
	ID          int64
	AppKey      string
	BotId       string
	Nickname    string
	BotPortrait string
	Description string
	BotType     int
	BotConf     string
	Status      BotStatus
}

type IBotConfStorage interface {
	Upsert(item BotConf) error
	FindById(appkey, botId string) (*BotConf, error)
	QryBotConfs(appkey string, startId, limit int64) ([]*BotConf, error)
	QryBotConfsWithStatus(appkey string, status BotStatus, startId, limit int64) ([]*BotConf, error)
}
