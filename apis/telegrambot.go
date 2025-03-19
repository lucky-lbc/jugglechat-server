package apis

import (
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func TelegramBotAdd(ctx *gin.Context) {
	req := apimodels.TelegramBot{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.TelegramBotAdd(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func TelegramBotDel(ctx *gin.Context) {
	req := apimodels.TelegramBot{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.TelegramBotDel(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func TelegramBotBatchDel(ctx *gin.Context) {
	req := apimodels.TelegramBotIds{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.TelegramBotBatchDel(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func TelegramBotList(ctx *gin.Context) {
	offset := ctx.Query("offset")
	count := 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}

	code, resp := services.QryTelegramBots(services.ToCtx(ctx), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}
