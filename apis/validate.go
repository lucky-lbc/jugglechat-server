package apis

import (
	"strings"

	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/gin-gonic/gin"
)

const (
	Header_RequestId     string = "request-id"
	Header_AppKey        string = "appkey"
	Header_Authorization string = "Authorization"
)

func Validate(ctx *gin.Context) {
	session := utils.GenerateUUIDShort11()
	ctx.Header(Header_RequestId, session)
	ctx.Set(string(services.CtxKey_Session), session)

	//check appkey
	appkey := ctx.Request.Header.Get(Header_AppKey)
	if appkey == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_APPKEY_REQUIRED)
		return
	}
	ctx.Set(string(services.CtxKey_AppKey), appkey)
	//check app exist
	appInfo, exist := services.GetAppInfo(appkey)
	if !exist {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
		ctx.Abort()
		return
	}
	urlPath := ctx.Request.URL.Path
	if urlPath != "/jim/login" && urlPath != "/jim/sms/send" && urlPath != "/jim/sms_login" && urlPath != "/jim/sms/login" && urlPath != "/jim/email/send" && urlPath != "/jim/email/login" && urlPath != "/jim/login/qrcode" && urlPath != "/jim/login/qrcode/check" {
		//current userId
		tokenStr := ctx.Request.Header.Get(Header_Authorization)
		if tokenStr == "" {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
			ctx.Abort()
			return
		}
		if tokenStr != "" {
			if strings.HasPrefix(tokenStr, "Bearer ") {
				tokenStr = tokenStr[7:]
				if !services.CheckApiKey(tokenStr, appkey, appInfo.AppSecureKey) {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					ctx.Abort()
					return
				}
			} else {
				authToken, err := services.ParseTokenString(tokenStr)
				if err != nil {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					ctx.Abort()
					return
				}
				token, err := services.ParseToken(authToken, []byte(appInfo.AppSecureKey))
				if err != nil {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					ctx.Abort()
					return
				}
				ctx.Set(string(services.CtxKey_RequesterId), token.UserId)
			}
		}
	}
}
