package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/commons/ctxs"
	"github.com/lucky-lbc/commons/errs"
	"github.com/lucky-lbc/commons/responses"
	"github.com/lucky-lbc/commons/tools"
	"github.com/lucky-lbc/jugglechat-server/admins/apis/models"
	"github.com/lucky-lbc/jugglechat-server/admins/services"
)

func QryUsers(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		responses.AdminErrorHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	offset := ctx.Query("offset")
	var count int64 = 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = tools.String2Int64(countStr)
		if err != nil {
			count = 20
		}
	}
	isPositiveOrder := false
	orderStr := ctx.Query("order")
	if orderStr != "" {
		order, err := strconv.Atoi(orderStr)
		if err == nil && order > 0 { //0:倒序;1:正序;
			isPositiveOrder = true
		}
	}
	code, users := services.QryUsers(ctxs.ToCtx(ctx), appkey, offset, count, isPositiveOrder)
	if code != errs.AdminErrorCode_Success {
		responses.AdminErrorHttpResp(ctx, code)
		return
	}
	responses.AdminSuccessHttpResp(ctx, users)
}

func BanUsers(ctx *gin.Context) {
	var req models.BanUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		responses.AdminErrorHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.BanUsers(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		responses.AdminErrorHttpResp(ctx, code)
		return
	}
	responses.AdminSuccessHttpResp(ctx, nil)
}

func UnBanUsers(ctx *gin.Context) {
	var req models.BanUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		responses.AdminErrorHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.UnBanUsers(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		responses.AdminErrorHttpResp(ctx, code)
		return
	}
	responses.AdminSuccessHttpResp(ctx, nil)
}
