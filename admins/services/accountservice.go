package services

import (
	"time"

	utils "github.com/juggleim/commons/tools"
	apimodels "github.com/juggleim/jugglechat-server/admins/apis/models"
	"github.com/juggleim/jugglechat-server/admins/errs"
	"github.com/juggleim/jugglechat-server/storages/dbs"
)

func CheckLogin(account, password string) (errs.AdminErrorCode, *apimodels.Account) {
	dao := dbs.AccountDao{}
	defaultAccount, err := dao.FindByAccount("admin")
	if err != nil || defaultAccount == nil {
		//init default account
		dao.Create(dbs.AccountDao{
			Account:       "admin",
			Password:      utils.SHA1("123456"),
			CreatedTime:   time.Now(),
			UpdatedTime:   time.Now(),
			State:         0,
			ParentAccount: "",
		})
	}

	password = utils.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err == nil && admin != nil {
		if admin.State != 0 {
			return errs.AdminErrorCode_AccountForbidden, nil
		}
		return errs.AdminErrorCode_Success, &apimodels.Account{
			Account:       admin.Account,
			State:         admin.State,
			ParentAccount: admin.ParentAccount,
			RoleId:        admin.RoleId,
			CreatedTime:   admin.CreatedTime.UnixMilli(),
			UpdatedTime:   admin.UpdatedTime.UnixMilli(),
		}
	}
	return errs.AdminErrorCode_LoginFail, nil
}

func CheckAccountState(account string) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	admin, err := dao.FindByAccount(account)
	if err != nil || admin == nil {
		return errs.AdminErrorCode_AccountNotExist
	} else {
		if admin.State == 0 {
			return errs.AdminErrorCode_Success
		} else {
			return errs.AdminErrorCode_AccountForbidden
		}
	}
}

func UpdPassword(account, password, newPassword string) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	password = utils.SHA1(password)
	admin, err := dao.FindByAccountPassword(account, password)
	if err != nil || admin == nil {
		return errs.AdminErrorCode_UpdPwdFail
	}
	newPassword = utils.SHA1(newPassword)
	dao.UpdatePassword(account, newPassword)
	return errs.AdminErrorCode_Success
}

func AddAccount(parentAccount, account, password string, roleId int) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	password = utils.SHA1(password)
	err := dao.Create(dbs.AccountDao{
		Account:       account,
		Password:      password,
		ParentAccount: parentAccount,
		UpdatedTime:   time.Now(),
		CreatedTime:   time.Now(),
		RoleId:        roleId,
	})
	if err != nil {
		return errs.AdminErrorCode_AccountExisted
	}
	return errs.AdminErrorCode_Success
}

func DisableAccounts(accounts []string, isDisable int) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	dao.UpdateState(accounts, isDisable)
	return errs.AdminErrorCode_Success
}
func DeleteAccounts(accounts []string) errs.AdminErrorCode {
	dao := dbs.AccountDao{}
	dao.DeleteAccounts(accounts)
	return errs.AdminErrorCode_Success
}

func QryAccounts(limit int64, offset string) *apimodels.Accounts {
	accounts := &apimodels.Accounts{
		Items:   []*apimodels.Account{},
		HasMore: false,
		Offset:  "",
	}
	dao := dbs.AccountDao{}
	offsetInt, err := utils.DecodeInt(offset)
	if err != nil {
		offsetInt = 0
	}
	dbAccounts, err := dao.QryAccounts(limit+1, offsetInt)
	if err == nil {
		if len(dbAccounts) > int(limit) {
			dbAccounts = dbAccounts[:len(dbAccounts)-1]
			accounts.HasMore = true
		}
		var id int64 = 0
		for _, dbAccount := range dbAccounts {
			accounts.Items = append(accounts.Items, &apimodels.Account{
				Account:       dbAccount.Account,
				State:         dbAccount.State,
				CreatedTime:   dbAccount.CreatedTime.UnixMilli(),
				UpdatedTime:   dbAccount.UpdatedTime.UnixMilli(),
				ParentAccount: dbAccount.ParentAccount,
				RoleId:        dbAccount.RoleId,
			})
			if dbAccount.ID > id {
				id = dbAccount.ID
			}
		}
		if id > 0 {
			offset, _ := utils.EncodeInt(id)
			accounts.Offset = offset
		}
	}
	return accounts
}
