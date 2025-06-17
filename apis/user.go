package apis

import (
	"bytes"
	"encoding/base64"
	"image/png"

	utils "github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/apis/responses"
	"github.com/juggleim/jugglechat-server/ctxs"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

func QryUserInfo(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, user := services.QryUserInfo(ctxs.ToCtx(ctx), userId)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, user)
}

func UpdateUser(ctx *gin.Context) {
	req := &models.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	services.UpdateUser(ctxs.ToCtx(ctx), req)
	responses.SuccessHttpResp(ctx, nil)
}

func UpdateUserSettings(ctx *gin.Context) {
	req := &models.UserSettings{}
	if err := ctx.BindJSON(req); err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateUserSettings(ctxs.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}

func SearchByPhone(ctx *gin.Context) {
	req := &models.UserObj{}
	if err := ctx.BindJSON(req); err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, users := services.SearchByPhone(ctxs.ToCtx(ctx), req.Phone)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, users)
}

func QryUserQrCode(ctx *gin.Context) {
	userId := ctx.GetString(string(ctxs.CtxKey_RequesterId))

	m := map[string]interface{}{
		"action":  "add_friend",
		"user_id": userId,
	}
	buf := bytes.NewBuffer([]byte{})
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	err := png.Encode(buf, qrCode)
	if err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	responses.SuccessHttpResp(ctx, map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}
