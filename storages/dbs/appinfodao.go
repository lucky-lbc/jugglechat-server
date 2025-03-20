package dbs

import (
	"fmt"
	"time"

	"github.com/juggleim/jugglechat-server/storages/dbs/dbcommons"
	"github.com/juggleim/jugglechat-server/storages/models"
)

type AppInfoDao struct {
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

func (app AppInfoDao) TableName() string {
	return "apps"
}

func (app AppInfoDao) Create(item models.AppInfo) error {
	err := dbcommons.GetDb().Create(&AppInfoDao{
		AppName:      item.AppName,
		AppKey:       item.AppKey,
		AppSecret:    item.AppSecret,
		AppSecureKey: item.AppSecureKey,
		AppStatus:    item.AppStatus,
		AppType:      item.AppType,
		CreatedTime:  item.CreatedTime,
		UpdatedTime:  item.UpdatedTime,
	}).Error
	return err
}

func (app AppInfoDao) Upsert(item models.AppInfo) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_name,app_key,app_secret,app_secure_key,app_type)VALUES(?,?,?,?,?)", app.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppName, item.AppKey, item.AppSecret, item.AppSecureKey, item.AppType).Error
}

func (app AppInfoDao) FindByAppkey(appkey string) (*models.AppInfo, error) {
	var appItem AppInfoDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Take(&appItem).Error
	if err != nil {
		return nil, err
	}
	return &models.AppInfo{
		ID:           appItem.ID,
		AppName:      appItem.AppName,
		AppKey:       appItem.AppKey,
		AppSecret:    appItem.AppSecret,
		AppSecureKey: appItem.AppSecureKey,
		AppStatus:    appItem.AppStatus,
		AppType:      appItem.AppType,
		CreatedTime:  appItem.CreatedTime,
		UpdatedTime:  appItem.UpdatedTime,
	}, nil
}
