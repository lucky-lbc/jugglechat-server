package apis

import (
	"github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/gin-gonic/gin"
)

func PostAdd(ctx *gin.Context) {
	req := models.Post{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PostAdd(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func QryPosts(ctx *gin.Context) {
	var limit int64 = 20
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = utils.String2Int64(limitStr)
		if err != nil {
			limit = 20
		}
	}
	var start int64
	startTimeStr := ctx.Query("start")
	start, err = utils.String2Int64(startTimeStr)
	if err != nil {
		start = 0
	}
	var isPositive bool = false
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err == nil {
		if order == 1 {
			isPositive = true
		}
	}
	code, resp := services.QryPosts(services.ToCtx(ctx), start, limit, isPositive)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PostInfo(ctx *gin.Context) {
	postId := ctx.Query("post_id")
	if postId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.QryPostInfo(services.ToCtx(ctx), postId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func QryPostComments(ctx *gin.Context) {
	postId := ctx.Query("post_id")
	if postId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	var limit int64 = 20
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = utils.String2Int64(limitStr)
		if err != nil {
			limit = 20
		}
	}
	var start int64
	startTimeStr := ctx.Query("start")
	start, err = utils.String2Int64(startTimeStr)
	if err != nil {
		start = 0
	}
	var isPositive bool = false
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err == nil {
		if order == 1 {
			isPositive = true
		}
	}
	code, resp := services.QryPostComments(services.ToCtx(ctx), postId, start, limit, isPositive)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func PostCommentAdd(ctx *gin.Context) {
	req := models.PostComment{}
	if err := ctx.BindJSON(&req); err != nil || req.PostId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.PostCommentAdd(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}
