package dbs

import (
	"bytes"
	"fmt"

	"github.com/lucky-lbc/commons/dbcommons"
	"github.com/lucky-lbc/jugglechat-server/storages/models"
)

type FriendRelDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	FriendId string `gorm:"friend_id"`
	OrderTag string `gorm:"order_tag"`
	AppKey   string `gorm:"app_key"`
}

func (rel FriendRelDao) TableName() string {
	return "friendrels"
}

func (rel FriendRelDao) Upsert(item models.FriendRel) error {
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,user_id,friend_id,order_tag)VALUES(?,?,?,?)", rel.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.UserId, item.FriendId, item.OrderTag).Error
}

func (rel FriendRelDao) BatchUpsert(items []models.FriendRel) error {
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,user_id,friend_id,order_tag)VALUES", rel.TableName())
	buffer.WriteString(sql)
	length := len(items)
	params := []interface{}{}
	for i, item := range items {
		if i == length-1 {
			buffer.WriteString("(?,?,?,?)")
		} else {
			buffer.WriteString("(?,?,?,?),")
		}
		params = append(params, item.AppKey, item.UserId, item.FriendId, item.OrderTag)
	}
	return dbcommons.GetDb().Exec(buffer.String(), params...).Error
}

func (rel FriendRelDao) QueryFriendRels(appkey, userId string, startId, limit int64) ([]*models.FriendRel, error) {
	var items []*FriendRelDao
	params := []interface{}{}
	condition := "app_key=?"
	params = append(params, appkey)
	if userId != "" {
		condition = condition + " and user_id=?"
		params = append(params, userId)
	}
	condition = condition + " and id>?"
	params = append(params, startId)
	err := dbcommons.GetDb().Where(condition, params...).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendRel{}
	for _, rel := range items {
		ret = append(ret, &models.FriendRel{
			ID:       rel.ID,
			AppKey:   rel.AppKey,
			UserId:   rel.UserId,
			FriendId: rel.FriendId,
			OrderTag: rel.OrderTag,
		})
	}
	return ret, nil
}

func (rel FriendRelDao) QueryFriendRelsWithPage(appkey, userId string, orderTag string, page, size int64) ([]*models.User, error) {
	sql := fmt.Sprintf("select r.*, u.nickname,u.user_portrait,u.user_type,u.pinyin from %s as r left join %s as u on r.app_key=u.app_key and r.friend_id=u.user_id where r.app_key=? and r.user_id=?", rel.TableName(), UserDao{}.TableName())
	params := []interface{}{}
	params = append(params, appkey, userId)
	if orderTag != "" {
		sql = sql + " and u.pinyin>=?"
		params = append(params, orderTag)
	}

	var items []*FriendRelWithUser
	err := dbcommons.GetDb().Raw(sql, params...).Order("case when u.pinyin REGEXP '^[A-Za-z]' then 1 else 0 end desc, u.pinyin asc").Offset((page - 1) * size).Limit(size).Find(&items).Error
	// var items []*FriendRelDao
	// params := []interface{}{}
	// condition := "app_key=? and user_id=?"
	// params = append(params, appkey, userId)
	// if orderTag != "" {
	// 	condition = condition + " and order_tag>=?"
	// 	params = append(params, orderTag)
	// }
	// err := dbcommons.GetDb().Where(condition, params...).Order("order_tag asc").Offset((page - 1) * size).Limit(size).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.User{}
	for _, item := range items {
		ret = append(ret, &models.User{
			UserId:       item.FriendId,
			Nickname:     item.Nickname,
			UserPortrait: item.UserPortrait,
			UserType:     item.UserType,
			Pinyin:       item.Pinyin,
			AppKey:       item.AppKey,
		})
	}
	return ret, nil
}

type FriendRelWithUser struct {
	FriendRelDao
	Nickname     string `gorm:"nickname"`
	UserPortrait string `gorm:"user_portrait"`
	UserType     int    `gorm:"user_type"`
	Pinyin       string `gorm:"pinyin"`
}

func (rel FriendRelDao) SearchFriendsByName(appkey, userId string, nickname string, startId, limit int64) ([]*models.User, error) {
	sql := fmt.Sprintf("select r.*,u.nickname,u.user_portrait,u.user_type,u.pinyin from %s as r left join %s as u on r.app_key=u.app_key and r.friend_id=u.user_id where r.app_key=? and r.user_id=? and r.id>? and u.nickname like ?", rel.TableName(), UserDao{}.TableName())
	var items []*FriendRelWithUser
	err := dbcommons.GetDb().Raw(sql, appkey, userId, startId, "%"+nickname+"%").Order("r.id asc").Limit(limit).Find(&items).Error
	ret := []*models.User{}
	if err == nil {
		for _, item := range items {
			ret = append(ret, &models.User{
				ID:           item.ID,
				UserId:       item.FriendId,
				Nickname:     item.Nickname,
				UserPortrait: item.UserPortrait,
				UserType:     item.UserType,
				Pinyin:       item.Pinyin,
				AppKey:       item.AppKey,
			})
		}
	}
	return ret, err
}

func (rel FriendRelDao) BatchDelete(appkey, userId string, friendIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and friend_id in (?)", appkey, userId, friendIds).Delete(&FriendRelDao{}).Error
}

func (rel FriendRelDao) QueryFriendRelsByFriendIds(appkey, userId string, friendIds []string) ([]*models.FriendRel, error) {
	var items []*FriendRelDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and friend_id in (?)", appkey, userId, friendIds).Order("id asc").Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.FriendRel{}
	for _, rel := range items {
		ret = append(ret, &models.FriendRel{
			ID:       rel.ID,
			AppKey:   rel.AppKey,
			UserId:   rel.UserId,
			FriendId: rel.FriendId,
			OrderTag: rel.OrderTag,
		})
	}
	return ret, nil
}

func (rel FriendRelDao) UpdateOrderTag(appkey, userId, friendId string, orderTag string) error {
	return dbcommons.GetDb().Model(&FriendRelDao{}).Where("app_key=? and user_id=? and friend_id=?", appkey, userId, friendId).Update("order_tag", orderTag).Error
}
