package apis

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

func QryUserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	code, user := services.QryUserInfo(services.ToCtx(ctx), userId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, user)
}

func UpdateUser(ctx *gin.Context) {
	req := &models.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	services.UpdateUser(services.ToCtx(ctx), req)
	SuccessHttpResp(ctx, nil)
}

func UpdateUserSettings(ctx *gin.Context) {
	req := &models.UserSettings{}
	if err := ctx.BindJSON(req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateUserSettings(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func SearchByPhone(ctx *gin.Context) {
	req := &models.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, users := services.SearchByPhone(services.ToCtx(ctx), req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, users)
}

func QryUserQrCode(ctx *gin.Context) {
	userId := ctx.GetString(string(services.CtxKey_RequesterId))

	m := map[string]interface{}{
		"action":  "add_friend",
		"user_id": userId,
	}
	buf := bytes.NewBuffer([]byte{})
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	err := png.Encode(buf, qrCode)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	SuccessHttpResp(ctx, map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}
