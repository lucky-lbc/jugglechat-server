package services

import (
	"jugglechat-server/services/sms"
	"jugglechat-server/storages"
	"jugglechat-server/utils/caches"
	"sync"
	"time"
)

var appCache *caches.LruCache
var appLock *sync.RWMutex

func init() {
	appCache = caches.NewLruCacheWithAddReadTimeout(1000, nil, 5*time.Minute, 5*time.Minute)
	appLock = &sync.RWMutex{}
}

type AppInfo struct {
	AppName      string `gorm:"app_name"`
	AppKey       string `gorm:"app_key"`
	AppSecret    string `gorm:"app_secret"`
	AppSecureKey string `gorm:"app_secure_key"`
	AppStatus    int    `gorm:"app_status"`

	SmsEngine sms.ISmsEngine
}

func GetAppInfo(appkey string) (*AppInfo, bool) {
	if obj, exist := appCache.Get(appkey); exist {
		return obj.(*AppInfo), true
	} else {
		appLock.Lock()
		defer appLock.Unlock()
		if obj, exist := appCache.Get(appkey); exist {
			return obj.(*AppInfo), true
		} else {
			storage := storages.NewAppInfoStorage()
			app, err := storage.FindByAppkey(appkey)
			if err == nil && app != nil {
				info := &AppInfo{
					AppName:      app.AppName,
					AppKey:       app.AppKey,
					AppSecret:    app.AppSecret,
					AppSecureKey: app.AppSecureKey,
					AppStatus:    app.AppStatus,
				}
				appCache.Add(appkey, info)
				return info, true
			}
			return nil, false
		}
	}
}
