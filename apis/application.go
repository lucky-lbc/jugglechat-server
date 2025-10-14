package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/commons/ctxs"
	"github.com/lucky-lbc/commons/errs"
	"github.com/lucky-lbc/commons/responses"
	"github.com/lucky-lbc/commons/tools"
	"github.com/lucky-lbc/jugglechat-server/services"
)

func QryApplications(ctx *gin.Context) {
	var err error
	var page int64 = 1
	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = tools.String2Int64(pageStr)
		if err != nil {
			page = 1
		}
	}
	var size int64 = 20
	sizeStr := ctx.Query("size")
	if sizeStr != "" {
		size, err = tools.String2Int64(sizeStr)
		if err != nil {
			size = 20
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
	code, resp := services.QryApplications(ctxs.ToCtx(ctx), page, size, isPositiveOrder)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, resp)
}
