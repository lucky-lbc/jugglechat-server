package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/commons/ctxs"
	"github.com/lucky-lbc/commons/errs"
	"github.com/lucky-lbc/commons/responses"
	"github.com/lucky-lbc/jugglechat-server/apis/models"
	"github.com/lucky-lbc/jugglechat-server/services"
)

func RecallMsg(ctx *gin.Context) {
	req := &models.RecallMsgReq{}
	if err := ctx.BindJSON(req); err != nil || req.TargetId == "" || req.ChannelType != 2 || req.MsgId == "" {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.RecallMsg(ctxs.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}

func DelMsgs(ctx *gin.Context) {
	req := &models.DelHisMsgsReq{}
	if err := ctx.BindJSON(req); err != nil || req.TargetId == "" || req.ChannelType != 2 || len(req.Msgs) <= 0 {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelMsgs(ctxs.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, nil)
}
