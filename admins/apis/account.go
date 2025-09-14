package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/responses"
	utils "github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/admins/services"
	"github.com/juggleim/jugglechat-server/storages/dbs"
)

func Login(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, account := services.CheckLogin(req.Account, req.Password)
	if code == errs.AdminErrorCode_Success {
		authStr, err := generateAuthorization(req.Account, dbs.RoleType(req.RoleId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &errs.AdminApiErrorMsg{
				Code: errs.AdminErrorCode_Default,
				Msg:  "auth fail",
			})
			return
		}
		responses.AdminSuccessHttpResp(ctx, &LoginResp{
			Account:       req.Account,
			Authorization: authStr,
			Env:           "private", //public
			// RoleId:        account.RoleId,
			RoleType: account.RoleType,
		})
	} else {
		ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
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
	// RoleId        int    `json:"role_id"`
	RoleType int `json:"role_type"`
}

func AddAccount(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	roleType := GetAccountRoleType(ctx)
	if roleType != dbs.RoleType_SuperAdmin {
		ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_NotPermission,
		})
		return
	}
	code := services.AddAccount(GetLoginedAccount(ctx), req.Account, req.Password, req.RoleId)
	ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
		Code: code,
	})
}

func UpdPassword(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.UpdPassword(req.Account, req.Password, req.NewPassword)
	ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
		Code: code,
	})
}

func DisableAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	roleType := GetAccountRoleType(ctx)
	if roleType != dbs.RoleType_SuperAdmin {
		ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_NotPermission,
		})
		return
	}
	code := services.DisableAccounts(req.Accounts, req.IsDisable)
	ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
		Code: code,
	})
}

func DeleteAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	roleType := GetAccountRoleType(ctx)
	if roleType != dbs.RoleType_SuperAdmin {
		ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_NotPermission,
		})
		return
	}
	code := services.DeleteAccounts(req.Accounts)
	ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
		Code: code,
	})
}

type AccountsReq struct {
	Accounts  []string `json:"accounts"`
	IsDisable int      `json:"is_disable"`
}

func QryAccounts(ctx *gin.Context) {
	roleType := GetAccountRoleType(ctx)
	if roleType != dbs.RoleType_SuperAdmin {
		ctx.JSON(http.StatusOK, &errs.AdminApiErrorMsg{
			Code: errs.AdminErrorCode_NotPermission,
		})
		return
	}
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
