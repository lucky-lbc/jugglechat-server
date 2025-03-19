package apis

import (
	"fmt"
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/services"

	"github.com/gin-gonic/gin"
)

func GetFileCred(ctx *gin.Context) {
	fmt.Print("xxxxxx")
	req := apimodels.QryFileCredReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.GetFileCred(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}
