package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/jugglechat-server/admins/apis/responses"
	"github.com/juggleim/jugglechat-server/admins/errs"
	"github.com/juggleim/jugglechat-server/admins/services"
	"github.com/juggleim/jugglechat-server/utils"
)

func Login(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, account := services.CheckLogin(req.Account, req.Password)
	if code == errs.AdminErrorCode_Success {
		authStr, err := generateAuthorization(req.Account)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &errs.ApiErrorMsg{
				Code: errs.AdminErrorCode_Default,
				Msg:  "auth fail",
			})
			return
		}
		responses.SuccessHttpResp(ctx, &LoginResp{
			Account:       req.Account,
			Authorization: authStr,
			Env:           "private", //public
			RoleId:        account.RoleId,
		})
	} else {
		ctx.JSON(http.StatusOK, &errs.ApiErrorMsg{
			Code: code,
			Msg:  "login failed",
		})
	}
}

type AccountReq struct {
	Account     string `json:"account"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	RoleId      int    `json:"role_id"`
}

type LoginResp struct {
	Account       string `json:"account"`
	Authorization string `json:"authorization"`
	Env           string `json:"env"`
	RoleId        int    `json:"role_id"`
}

func AddAccount(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.AddAccount(GetLoginedAccount(ctx), req.Account, req.Password, req.RoleId)
	ctx.JSON(http.StatusOK, &errs.ApiErrorMsg{
		Code: code,
	})
}

func UpdPassword(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.UpdPassword(req.Account, req.Password, req.NewPassword)
	ctx.JSON(http.StatusOK, &errs.ApiErrorMsg{
		Code: code,
	})
}

func DisableAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &errs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.DisableAccounts(req.Accounts, req.IsDisable)
	ctx.JSON(http.StatusOK, &errs.ApiErrorMsg{
		Code: code,
	})
}

func DeleteAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &errs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.DeleteAccounts(req.Accounts)
	ctx.JSON(http.StatusOK, &errs.ApiErrorMsg{
		Code: code,
	})
}

type AccountsReq struct {
	Accounts  []string `json:"accounts"`
	IsDisable int      `json:"is_disable"`
}

func QryAccounts(ctx *gin.Context) {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := utils.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	accounts := services.QryAccounts(limit, offsetStr)
	responses.SuccessHttpResp(ctx, accounts)
}
