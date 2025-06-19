package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/admins/apis/responses"
	"github.com/juggleim/jugglechat-server/admins/errs"
	"github.com/juggleim/jugglechat-server/admins/services"
	"github.com/juggleim/jugglechat-server/ctxs"
)

func QryUsers(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		responses.ErrorHttpResp(ctx, errs.AdminErrorCode_ParamError)
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
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, users)
}
