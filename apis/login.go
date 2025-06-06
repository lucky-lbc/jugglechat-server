package apis

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/services/imsdk"
	"github.com/juggleim/jugglechat-server/storages"
	"github.com/juggleim/jugglechat-server/storages/models"
	"github.com/juggleim/jugglechat-server/utils"

	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func Login(ctx *gin.Context) {
	req := &apimodels.LoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Account == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := utils.ShortMd5(req.Account)
	nickname := fmt.Sprintf("user%05d", utils.RandInt(100000))
	fmt.Println(userId, nickname)
	appkey := ctx.GetString(string(services.CtxKey_AppKey))
	sdk := imsdk.GetImSdk(appkey)
	if sdk == nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
		return
	}
	resp, code, _, err := sdk.Register(juggleimsdk.User{
		UserId:   userId,
		Nickname: nickname,
	})
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
		return
	}
	if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
		ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	SuccessHttpResp(ctx, resp)
}

func SmsSend(ctx *gin.Context) {
	req := &apimodels.SmsLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Phone == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SmsSend(services.ToCtx(ctx), req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func SmsLogin(ctx *gin.Context) {
	req := &apimodels.SmsLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Phone == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.CheckPhoneSmsCode(services.ToCtx(ctx), req.Phone, req.Code)
	if code == errs.IMErrorCode_SUCCESS {
		appkey := ctx.GetString(string(services.CtxKey_AppKey))
		userId := utils.ShortMd5(req.Phone)
		nickname := fmt.Sprintf("user%05d", utils.RandInt(100000))
		storage := storages.NewUserStorage()
		user, err := storage.FindByPhone(appkey, req.Phone)
		if err == nil && user != nil {
			userId = user.UserId
			nickname = user.Nickname
		} else {
			user, err = storage.FindByUserId(appkey, userId)
			if err == nil && user != nil {
				userId = user.UserId
				nickname = user.Nickname
			} else {
				if err != gorm.ErrRecordNotFound {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				userId = utils.GenerateUUIDShort11()
				err = storage.Create(models.User{
					UserId:   userId,
					Nickname: nickname,
					Phone:    req.Phone,
					AppKey:   appkey,
				})
				if err != nil {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				} else {
					userExtStorage := storages.NewUserExtStorage()
					userExtStorage.Upsert(models.UserExt{
						UserId:    userId,
						ItemKey:   apimodels.UserExtKey_FriendVerifyType,
						ItemValue: utils.Int2String(int64(apimodels.FriendVerifyType_NeedFriendVerify)),
						ItemType:  apimodels.AttItemType_Setting,
						AppKey:    appkey,
					})
				}
				//assistant send welcome message
				services.InitUserAssistant(services.ToCtx(ctx), userId, nickname, "")
			}
		}
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: nickname,
		})
		if err != nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}

		SuccessHttpResp(ctx, &apimodels.LoginUserResp{
			UserId:        userId,
			NickName:      nickname,
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else {
		ErrorHttpResp(ctx, code)
		return
	}
}

func EmailSend(ctx *gin.Context) {
	req := &apimodels.EmailLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Email == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	// code := services.SmsSend(ctx.ToRpcCtx(), req.Email)
	// if code != errs.IMErrorCode_SUCCESS {
	// 	ctx.ResponseErr(code)
	// 	return
	// }
	SuccessHttpResp(ctx, nil)
}

func EmailLogin(ctx *gin.Context) {
	req := &apimodels.EmailLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Email == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.CheckEmailCode(services.ToCtx(ctx), req.Email, req.Code)
	if code == errs.IMErrorCode_SUCCESS {
		appkey := ctx.GetString(string(services.CtxKey_AppKey))
		userId := utils.ShortMd5(req.Email)
		nickname := fmt.Sprintf("user%05d", utils.RandInt(100000))
		storage := storages.NewUserStorage()
		user, err := storage.FindByPhone(appkey, req.Email)
		if err == nil && user != nil {
			userId = user.UserId
			nickname = user.Nickname
		} else {
			user, err = storage.FindByUserId(appkey, userId)
			if err == nil && user != nil {
				userId = user.UserId
				nickname = user.Nickname
			} else {
				if err != gorm.ErrRecordNotFound {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				userId = utils.GenerateUUIDShort11()
				err = storage.Create(models.User{
					UserId:   userId,
					Nickname: nickname,
					Email:    req.Email,
					AppKey:   appkey,
				})
				if err != nil {
					ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				} else {
					userExtStorage := storages.NewUserExtStorage()
					userExtStorage.Upsert(models.UserExt{
						UserId:    userId,
						ItemKey:   apimodels.UserExtKey_FriendVerifyType,
						ItemValue: utils.Int2String(int64(apimodels.FriendVerifyType_NeedFriendVerify)),
						ItemType:  apimodels.AttItemType_Setting,
						AppKey:    appkey,
					})
				}
				//assistant send welcome message
				services.InitUserAssistant(services.ToCtx(ctx), userId, nickname, "")
			}
		}
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: nickname,
		})
		if err != nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}

		SuccessHttpResp(ctx, &apimodels.LoginUserResp{
			UserId:        userId,
			NickName:      nickname,
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else {
		ErrorHttpResp(ctx, code)
		return
	}
}

func GenerateQrCode(ctx *gin.Context) {
	uuidStr := utils.GenerateUUIDString()
	m := map[string]interface{}{
		"action": "login",
		"code":   uuidStr,
	}
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, qrCode)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	storage.Create(models.QrCodeRecord{
		CodeId:      uuidStr,
		AppKey:      ctx.GetString(string(services.CtxKey_AppKey)),
		CreatedTime: time.Now().UnixMilli(),
	})
	SuccessHttpResp(ctx, map[string]string{
		"id":      uuidStr,
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}

func CheckQrCode(ctx *gin.Context) {
	req := &apimodels.QrCode{}
	if err := ctx.BindJSON(req); err != nil || req.Id == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	record, err := storage.FindById(ctx.GetString(string(services.CtxKey_AppKey)), req.Id)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	if time.Now().UnixMilli()-record.CreatedTime > 10*60*1000 {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_QRCODE_EXPIRED)
		return
	}
	appkey := ctx.GetString(string(services.CtxKey_AppKey))
	if record.Status == models.QrCodeRecordStatus_OK {
		userId := record.UserId
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: "",
		})
		if err != nil {
			ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}
		SuccessHttpResp(ctx, &apimodels.LoginUserResp{
			UserId:        userId,
			NickName:      "",
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else if record.Status == models.QrCodeRecordStatus_Default {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_CONTINUE)
		return
	}
}

func ConfirmQrCode(ctx *gin.Context) {
	req := &apimodels.QrCode{}
	if err := ctx.BindJSON(req); err != nil || req.Id == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	appkey := ctx.GetString(string(services.CtxKey_AppKey))
	userId := ctx.GetString(string(services.CtxKey_RequesterId))
	storage := storages.NewQrCodeRecordStorage()
	err := storage.UpdateStatus(appkey, req.Id, models.QrCodeRecordStatus_OK, userId)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	SuccessHttpResp(ctx, nil)
}
