package apis

import (
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/services"

	"github.com/gin-gonic/gin"
)

func ChgGroupOwner(ctx *gin.Context) {
	req := &apimodels.GroupOwnerChgReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ChgGroupOwner(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func AddGrpAdministrator(ctx *gin.Context) {
	req := &apimodels.GroupAdministratorsReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddGroupAdministrators(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func DelGrpAdministrator(ctx *gin.Context) {
	req := &apimodels.GroupAdministratorsReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGroupAdministrators(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QryGrpAdministrators(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	code, resp := services.QryGroupAdministrators(services.ToCtx(ctx), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func SetGroupMute(ctx *gin.Context) {
	req := &apimodels.SetGroupMuteReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupMute(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func SetGrpVerifyType(ctx *gin.Context) {
	req := &apimodels.SetGroupVerifyTypeReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupVerifyType(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func SetGrpHisMsgVisible(ctx *gin.Context) {
	req := &apimodels.SetGroupHisMsgVisibleReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGroupHisMsgVisible(services.ToCtx(ctx), req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}
