package services

import (
	"context"

	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/tools"
	apimodels "github.com/juggleim/jugglechat-server/admins/apis/models"
	"github.com/juggleim/jugglechat-server/storages"
)

func QryUsers(ctx context.Context, appkey, offset string, limit int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Users) {
	var startId int64 = 0
	var err error
	if offset != "" {
		startId, err = tools.DecodeInt(offset)
		if err != nil {
			startId = 0
		}
	}
	ret := &apimodels.Users{
		Items: []*apimodels.User{},
	}
	storage := storages.NewUserStorage()
	users, err := storage.QryUsers(appkey, startId, limit, isPositive)
	if err == nil {
		for _, user := range users {
			ret.Offset, _ = tools.EncodeInt(user.ID)
			ret.Items = append(ret.Items, &apimodels.User{
				UserId:      user.UserId,
				Nickname:    user.Nickname,
				Avatar:      user.UserPortrait,
				Pinyin:      user.Pinyin,
				UserType:    user.UserType,
				Status:      int32(user.Status),
				CreatedTime: user.CreatedTime.UnixMilli(),
			})
		}
	}
	return errs.AdminErrorCode_Success, ret
}

func BanUsers(ctx context.Context) {

}
