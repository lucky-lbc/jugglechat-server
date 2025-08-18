package models

type Application struct {
	ID          int64
	AppId       string
	AppName     string
	AppIcon     string
	AppDesc     string
	AppUrl      string
	AppOrder    int
	CreatedTime int64
	UpdatedTime int64
	AppKey      string
}

type IApplicationStorage interface {
	Create(item Application) error
	FindByAppId(appkey, appId string) (*Application, error)
	QryApplications(appkey string, limit int64) ([]*Application, error)
}
