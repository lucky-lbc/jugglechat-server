package imsdk

import (
	"sync"

	"github.com/juggleim/commons/configures"
	"github.com/juggleim/commons/dbcommons"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

var imsdkMap *sync.Map
var imLock *sync.RWMutex

func init() {
	imsdkMap = &sync.Map{}
	imLock = &sync.RWMutex{}
}

func GetImSdk(appkey string) *juggleimsdk.JuggleIMSdk {
	if val, exist := imsdkMap.Load(appkey); exist {
		return val.(*juggleimsdk.JuggleIMSdk)
	} else {
		imLock.Lock()
		defer imLock.Unlock()

		if val, exist := imsdkMap.Load(appkey); exist {
			return val.(*juggleimsdk.JuggleIMSdk)
		} else {
			dao := dbcommons.AppInfoDao{}
			appinfo, _ := dao.FindByAppkey(appkey)
			if appinfo != nil {
				sdk := juggleimsdk.NewJuggleIMSdk(appkey, appinfo.AppSecret, configures.Config.ImApiDomain)
				imsdkMap.Store(appkey, sdk)
				return sdk
			}
			return nil
		}
	}
}
