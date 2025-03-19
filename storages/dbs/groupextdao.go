package dbs

import (
	"fmt"
	"jugglechat-server/storages/dbs/dbcommons"
	"jugglechat-server/storages/models"
	"time"
)

type GroupExtDao struct {
	ID          int64     `gorm:"primary_key"`
	GroupId     string    `gorm:"group_id"`
	ItemKey     string    `gorm:"item_key"`
	ItemValue   string    `gorm:"item_value"`
	ItemType    int       `gorm:"item_type"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (ext GroupExtDao) TableName() string {
	return "groupinfoexts"
}

func (ext GroupExtDao) Find(appkey, groupId string, itemKey string) (*models.GroupExt, error) {
	var item GroupExtDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=? and app_item_key=?", appkey, groupId, itemKey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.GroupExt{
		ID:          item.ID,
		GroupId:     item.GroupId,
		ItemKey:     item.ItemKey,
		ItemValue:   item.ItemValue,
		ItemType:    item.ItemType,
		UpdatedTime: item.UpdatedTime,
		AppKey:      item.AppKey,
	}, nil
}

func (ext GroupExtDao) QryExtFields(appkey, groupId string) ([]*models.GroupExt, error) {
	var items []*GroupExtDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Find(&items).Error
	ret := []*models.GroupExt{}
	for _, item := range items {
		ret = append(ret, &models.GroupExt{
			ID:          item.ID,
			GroupId:     item.GroupId,
			ItemKey:     item.ItemKey,
			ItemValue:   item.ItemValue,
			ItemType:    item.ItemType,
			UpdatedTime: item.UpdatedTime,
			AppKey:      item.AppKey,
		})
	}
	return ret, err
}

func (ext GroupExtDao) Upsert(item models.GroupExt) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_key,group_id,item_key,item_value,item_type)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE item_value=?", ext.TableName()), item.AppKey, item.GroupId, item.ItemKey, item.ItemValue, item.ItemType, item.ItemValue).Error
}
