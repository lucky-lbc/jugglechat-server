package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/jugglechat-server/apis/models"
	"github.com/lucky-lbc/jugglechat-server/commons/ctxs"
	"github.com/lucky-lbc/jugglechat-server/commons/errs"
	"github.com/lucky-lbc/jugglechat-server/commons/responses"
	"github.com/lucky-lbc/jugglechat-server/services"
)

func AddFeedback(ctx *gin.Context) {
	req := models.Feedback{}
	if err := ctx.BindJSON(&req); err != nil || req.Category == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddFeedback(ctxs.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}
