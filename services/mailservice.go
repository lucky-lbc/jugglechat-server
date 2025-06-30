package services

import (
	"context"
	"fmt"
	"time"

	"github.com/juggleim/commons/appinfos"
	"github.com/juggleim/commons/ctxs"
	"github.com/juggleim/commons/dbcommons"
	"github.com/juggleim/commons/emailengines"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/storages"
	"github.com/juggleim/jugglechat-server/storages/models"
)

func MailSend(ctx context.Context, mailAddress string) errs.IMErrorCode {
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	mailEngine := GetMailEngine(appkey)
	if mailEngine != nil {
		//检查是否还有有效的
		storage := storages.NewSmsRecordStorage()
		record, err := storage.FindByEmail(appkey, mailAddress, time.Now().Add(-3*time.Minute))
		randomCode := RandomSms()
		if err == nil {
			randomCode = record.Code
		} else {
			_, err = storage.Create(models.SmsRecord{
				AppKey:      appkey,
				Email:       mailAddress,
				Code:        randomCode,
				CreatedTime: time.Now(),
			})
			if err != nil {
				return errs.IMErrorCode_APP_SMS_SEND_FAILED
			}
		}
		err = mailEngine.SendMail(mailAddress, "Login Verify", fmt.Sprintf("Your login verification code is: %s", randomCode))
		if err != nil {
			return errs.IMErrorCode_APP_SMS_SEND_FAILED
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func CheckEmailCode(ctx context.Context, mail, code string) errs.IMErrorCode {
	if code == "000000" {
		return errs.IMErrorCode_SUCCESS
	}
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	storage := storages.NewSmsRecordStorage()
	record, err := storage.FindByEmailCode(appkey, mail, code)
	if err != nil {
		return errs.IMErrorCode_APP_SMS_CODE_EXPIRED
	}
	interval := time.Since(record.CreatedTime)
	if interval > 5*time.Minute {
		return errs.IMErrorCode_APP_SMS_CODE_EXPIRED
	}
	return errs.IMErrorCode_SUCCESS
}

func GetMailEngine(appkey string) emailengines.IEmailEngine {
	appInfo, exist := appinfos.GetAppInfo(appkey)
	if exist && appInfo != nil {
		if appInfo.MailEngine == nil {
			appinfos.GetAppLock().Lock()
			defer appinfos.GetAppLock().Unlock()
			loadMailEngine(appInfo)
		}
		if appInfo.MailEngine != nil {
			return appInfo.MailEngine
		}
	}
	return nil
}

func loadMailEngine(appInfo *appinfos.AppInfo) {
	extDao := dbcommons.AppExtDao{}
	ext, err := extDao.Find(appInfo.AppKey, "mail_engine_conf")
	if err == nil && ext.AppItemValue != "" {
		mailConf := &MailEngineConf{}
		err := tools.JsonUnMarshal([]byte(ext.AppItemValue), mailConf)
		if err == nil {
			if mailConf.Channel == "ali" && mailConf.AliMailEngine != nil && mailConf.AliMailEngine.AccessKeyId != "" && mailConf.AliMailEngine.AccessKeySecret != "" {
				appInfo.MailEngine = mailConf.AliMailEngine
				return
			}
		}
	}
	appInfo.MailEngine = nil
}

type MailEngineConf struct {
	Channel       string                       `json:"channel,omitempty"`
	AliMailEngine *emailengines.AliEmailEngine `json:"ali,omitempty"`
}
