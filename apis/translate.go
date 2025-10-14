package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/commons/ctxs"
	"github.com/lucky-lbc/commons/errs"
	"github.com/lucky-lbc/commons/responses"
	"github.com/lucky-lbc/jugglechat-server/apis/models"
	"github.com/lucky-lbc/jugglechat-server/services"
)

func Translate(ctx *gin.Context) {
	req := &models.TransReq{}
	if err := ctx.BindJSON(req); err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.Translate(ctxs.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, resp)
}
