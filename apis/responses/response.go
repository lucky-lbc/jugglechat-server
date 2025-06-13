package responses

import (
	"net/http"

	"github.com/juggleim/jugglechat-server/errs"

	"github.com/gin-gonic/gin"
)

func ErrorHttpResp(ctx *gin.Context, code errs.IMErrorCode) {
	apiErr := errs.GetApiErrorByCode(code)
	ctx.JSON(apiErr.HttpCode, apiErr)
}

func SuccessHttpResp(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, errs.SuccHttpResp{
		ApiErrorMsg: errs.ApiErrorMsg{
			Code: 0,
			Msg:  "success",
		},
		Data: data,
	})
}
