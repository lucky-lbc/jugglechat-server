package apis

import (
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AssistantAnswer(ctx *gin.Context) {
	req := apimodels.AssistantAnswerReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.AutoAnswer(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PromptAdd(ctx *gin.Context) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Prompts == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PromptAdd(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PromptUpdate(ctx *gin.Context) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Id == "" || req.Prompts == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptUpdate(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func PromptDel(ctx *gin.Context) {
	req := apimodels.Prompt{}
	if err := ctx.BindJSON(&req); err != nil || req.Id == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptDel(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func PromptBatchDel(ctx *gin.Context) {
	req := apimodels.PromptIds{}
	if err := ctx.BindJSON(&req); err != nil || len(req.Ids) <= 0 {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.PromptBatchDel(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QryPrompts(ctx *gin.Context) {
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
	code, prompts := services.QryPrompts(services.ToCtx(ctx), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, prompts)
}
