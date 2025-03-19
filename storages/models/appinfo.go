package models

import "time"

type AppStatus int
type AppType int

var (
	AppStatus_Normal AppStatus = 0
	AppStatus_Block  AppStatus = 1
	AppStatus_Expire AppStatus = 2

	AppType_Private AppType = 0
	AppType_Alone   AppType = 1
	AppType_Public  AppType = 2
)

type AppInfo struct {
	ID           int64     `gorm:"primary_key"`
	AppName      string    `gorm:"app_name"`
	AppKey       string    `gorm:"app_key"`
	AppSecret    string    `gorm:"app_secret"`
	AppSecureKey string    `gorm:"app_secure_key"`
	AppStatus    int       `gorm:"app_status"`
	AppType      int       `gorm:"app_type"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdatedTime  time.Time `gorm:"updated_time"`
}

type IAppInfoStorage interface {
	Create(item AppInfo) error
	Upsert(item AppInfo) error
	FindByAppkey(appkey string) (*AppInfo, error)
}

type AppExt struct {
	AppKey       string
	AppItemKey   string
	AppItemValue string
	UpdatedTime  time.Time
}

type IAppExtStorage interface {
	FindListByAppkey(appkey string) ([]*AppExt, error)
	Find(appkey string, itemKey string) (*AppExt, error)
	FindByItemKeys(appkey string, itemKeys []string) ([]*AppExt, error)
	Upsert(appkey string, fieldKey, fieldValue string) error
}
