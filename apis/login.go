package apis

import (
	"bytes"
	"encoding/base64"
	"fmt"

	utils "github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/apis/responses"
	"github.com/juggleim/jugglechat-server/ctxs"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/events"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/services/imsdk"
	"github.com/juggleim/jugglechat-server/storages"
	dbModels "github.com/juggleim/jugglechat-server/storages/models"

	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func Login(ctx *gin.Context) {
	req := &models.LoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Account == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := utils.ShortMd5(req.Account)
	nickname := fmt.Sprintf("user%05d", utils.RandInt(100000))
	fmt.Println(userId, nickname)
	appkey := ctx.GetString(string(ctxs.CtxKey_AppKey))
	sdk := imsdk.GetImSdk(appkey)
	if sdk == nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
		return
	}
	resp, code, _, err := sdk.Register(juggleimsdk.User{
		UserId:   userId,
		Nickname: nickname,
	})
	if err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
		return
	}
	if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	responses.SuccessHttpResp(ctx, resp)
}

func SmsSend(ctx *gin.Context) {
	req := &models.SmsLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Phone == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SmsSend(ctxs.ToCtx(ctx), req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}

func SmsLogin(ctx *gin.Context) {
	req := &models.SmsLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Phone == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.CheckPhoneSmsCode(ctxs.ToCtx(ctx), req.Phone, req.Code)
	if code == errs.IMErrorCode_SUCCESS {
		appkey := ctx.GetString(string(ctxs.CtxKey_AppKey))
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
					responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				}
				userId = utils.GenerateUUIDShort11()
				err = storage.Create(dbModels.User{
					UserId:   userId,
					Nickname: nickname,
					Phone:    req.Phone,
					AppKey:   appkey,
				})
				if err != nil {
					responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
					return
				} else {
					events.TriggerUserRegiste(dbModels.User{
						UserId:   userId,
						Nickname: nickname,
						Phone:    req.Phone,
						AppKey:   appkey,
					})
					userExtStorage := storages.NewUserExtStorage()
					userExtStorage.Upsert(dbModels.UserExt{
						UserId:    userId,
						ItemKey:   models.UserExtKey_FriendVerifyType,
						ItemValue: utils.Int2String(int64(models.FriendVerifyType_NeedFriendVerify)),
						ItemType:  models.AttItemType_Setting,
						AppKey:    appkey,
					})
				}
			}
		}
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: nickname,
		})
		if err != nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}

		responses.SuccessHttpResp(ctx, &models.LoginUserResp{
			UserId:        userId,
			NickName:      nickname,
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else {
		responses.ErrorHttpResp(ctx, code)
		return
	}
}

func EmailSend(ctx *gin.Context) {
	req := &models.EmailLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Email == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	// code := services.SmsSend(ctx.ToRpcCtx(), req.Email)
	// if code != errs.IMErrorCode_SUCCESS {
	// 	ctx.ResponseErr(code)
	// 	return
	// }
	responses.SuccessHttpResp(ctx, nil)
}

func EmailLogin(ctx *gin.Context) {
	req := &models.EmailLoginReq{}
	if err := ctx.BindJSON(req); err != nil || req.Email == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.CheckEmailCode(ctxs.ToCtx(ctx), req.Email, req.Code)
	if code == errs.IMErrorCode_SUCCESS {
		appkey := ctx.GetString(string(ctxs.CtxKey_AppKey))
		var userId string
		nickname := fmt.Sprintf("user%05d", utils.RandInt(100000))
		storage := storages.NewUserStorage()
		user, err := storage.FindByEmail(appkey, req.Email)
		if err == nil && user != nil {
			userId = user.UserId
			nickname = user.Nickname
		} else {
			if err != gorm.ErrRecordNotFound {
				responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
				return
			}
			userId = utils.GenerateUUIDShort11()
			err = storage.Create(dbModels.User{
				UserId:   userId,
				Nickname: nickname,
				Email:    req.Email,
				AppKey:   appkey,
			})
			if err != nil {
				responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_LOGIN)
				return
			} else {
				events.TriggerUserRegiste(dbModels.User{
					UserId:   userId,
					Nickname: nickname,
					Email:    req.Email,
					AppKey:   appkey,
				})
				userExtStorage := storages.NewUserExtStorage()
				userExtStorage.Upsert(dbModels.UserExt{
					UserId:    userId,
					ItemKey:   models.UserExtKey_FriendVerifyType,
					ItemValue: utils.Int2String(int64(models.FriendVerifyType_NeedFriendVerify)),
					ItemType:  models.AttItemType_Setting,
					AppKey:    appkey,
				})
			}
		}
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: nickname,
		})
		if err != nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}

		responses.SuccessHttpResp(ctx, &models.LoginUserResp{
			UserId:        userId,
			NickName:      nickname,
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else {
		responses.ErrorHttpResp(ctx, code)
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
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	storage.Create(dbModels.QrCodeRecord{
		CodeId:      uuidStr,
		AppKey:      ctx.GetString(string(ctxs.CtxKey_AppKey)),
		CreatedTime: time.Now().UnixMilli(),
	})
	responses.SuccessHttpResp(ctx, map[string]string{
		"id":      uuidStr,
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}

func CheckQrCode(ctx *gin.Context) {
	req := &models.QrCode{}
	if err := ctx.BindJSON(req); err != nil || req.Id == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	storage := storages.NewQrCodeRecordStorage()
	record, err := storage.FindById(ctx.GetString(string(ctxs.CtxKey_AppKey)), req.Id)
	if err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	if time.Now().UnixMilli()-record.CreatedTime > 10*60*1000 {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_QRCODE_EXPIRED)
		return
	}
	appkey := ctx.GetString(string(ctxs.CtxKey_AppKey))
	if record.Status == dbModels.QrCodeRecordStatus_OK {
		userId := record.UserId
		sdk := imsdk.GetImSdk(appkey)
		if sdk == nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_NOT_EXISTED)
			return
		}
		resp, code, _, err := sdk.Register(juggleimsdk.User{
			UserId:   userId,
			Nickname: "",
		})
		if err != nil {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_INTERNAL_TIMEOUT)
			return
		}
		if code != juggleimsdk.ApiCode(errs.IMErrorCode_SUCCESS) {
			responses.ErrorHttpResp(ctx, errs.IMErrorCode(code))
			return
		}
		responses.SuccessHttpResp(ctx, &models.LoginUserResp{
			UserId:        userId,
			NickName:      "",
			Authorization: services.GenerateToken(appkey, userId),
			ImToken:       resp.Token,
		})
	} else if record.Status == dbModels.QrCodeRecordStatus_Default {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_CONTINUE)
		return
	}
}

func ConfirmQrCode(ctx *gin.Context) {
	req := &models.QrCode{}
	if err := ctx.BindJSON(req); err != nil || req.Id == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	appkey := ctx.GetString(string(ctxs.CtxKey_AppKey))
	userId := ctx.GetString(string(ctxs.CtxKey_RequesterId))
	storage := storages.NewQrCodeRecordStorage()
	err := storage.UpdateStatus(appkey, req.Id, dbModels.QrCodeRecordStatus_OK, userId)
	if err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}
