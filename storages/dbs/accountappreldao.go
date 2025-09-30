package dbs

import (
	"fmt"

	"github.com/juggleim/commons/dbcommons"
)

type AccountAppRelDao struct {
	ID      int64  `gorm:"primary_key"`
	AppKey  string `gorm:"app_key"`
	Account string `gorm:"account"`
}

func (rel AccountAppRelDao) TableName() string {
	return "accountapprels"
}

func (rel AccountAppRelDao) Create(item AccountAppRelDao) error {
	return dbcommons.GetDb().Create(&item).Error
}

func (rel AccountAppRelDao) CheckExist(appkey string, account string) bool {
	var item AccountAppRelDao
	err := dbcommons.GetDb().Where("account=? and app_key=?", account, appkey).Take(&item).Error
	return err == nil
}

func (rel AccountAppRelDao) BatchDelete(account string, appkeys []string) error {
	return dbcommons.GetDb().Where("account=? and app_key in (?)", account, appkeys).Delete(&rel).Error
}

func (rel AccountAppRelDao) FindByAppkey(account string, appkey string) *dbcommons.AppInfoDao {
	var appItem dbcommons.AppInfoDao
	sql := fmt.Sprintf("select app.* from %s as rel left join %s as app on rel.app_key=app.app_key where rel.account=? and rel.app_key=?", rel.TableName(), dbcommons.AppInfoDao{}.TableName())
	err := dbcommons.GetDb().Raw(sql, account, appkey).Take(&appItem).Error
	if err != nil {
		return nil
	}
	return &appItem
}

func (rel AccountAppRelDao) QryApps(account string, limit int64, offset int64) ([]*dbcommons.AppInfoDao, error) {
	var list []*dbcommons.AppInfoDao
	sql := fmt.Sprintf("select app.* from %s as rel left join %s as app on rel.app_key=app.app_key where rel.account=? and app.id<?", rel.TableName(), dbcommons.AppInfoDao{}.TableName())
	err := dbcommons.GetDb().Raw(sql, account, offset).Order("app.id desc").Limit(limit).Find(&list).Error
	return list, err
}
