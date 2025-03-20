package dbs

import (
	"fmt"
	"time"

	"github.com/juggleim/jugglechat-server/storages/dbs/dbcommons"
	"github.com/juggleim/jugglechat-server/storages/models"
)

type AppExtDao struct {
	ID           int64     `gorm:"primary_key"`
	AppKey       string    `gorm:"app_key"`
	AppItemKey   string    `gorm:"app_item_key"`
	AppItemValue string    `gorm:"app_item_value"`
	UpdatedTime  time.Time `gorm:"updated_time"`
}

func (appExt AppExtDao) TableName() string {
	return "appexts"
}

func (appExt AppExtDao) FindListByAppkey(appkey string) ([]*models.AppExt, error) {
	var list []*AppExtDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Find(&list).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.AppExt{}
	for _, item := range list {
		ret = append(ret, &models.AppExt{
			AppKey:       item.AppKey,
			AppItemKey:   item.AppItemKey,
			AppItemValue: item.AppItemValue,
			UpdatedTime:  item.UpdatedTime,
		})
	}
	return ret, nil
}

func (appExt AppExtDao) Find(appkey string, itemKey string) (*models.AppExt, error) {
	var item AppExtDao
	err := dbcommons.GetDb().Where("app_key=? and app_item_key=?", appkey, itemKey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.AppExt{
		AppKey:       item.AppKey,
		AppItemKey:   item.AppItemKey,
		AppItemValue: item.AppItemValue,
		UpdatedTime:  item.UpdatedTime,
	}, nil
}

func (appExt AppExtDao) FindByItemKeys(appkey string, itemKeys []string) ([]*models.AppExt, error) {
	var list []*AppExtDao
	err := dbcommons.GetDb().Where("app_key=? and app_item_key in(?)", appkey, itemKeys).Find(&list).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.AppExt{}
	for _, item := range list {
		ret = append(ret, &models.AppExt{
			AppKey:       item.AppKey,
			AppItemKey:   item.AppItemKey,
			AppItemValue: item.AppItemValue,
			UpdatedTime:  item.UpdatedTime,
		})
	}
	return ret, err
}

func (appExt AppExtDao) Upsert(appkey string, fieldKey, fieldValue string) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,app_item_key,app_item_value)VALUES(?,?,?) ON DUPLICATE KEY UPDATE app_item_value=?", appExt.TableName()), appkey, fieldKey, fieldValue, fieldValue).Error
}
