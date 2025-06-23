package services

import (
	"context"

	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/tools"
	apimodels "github.com/juggleim/jugglechat-server/admins/apis/models"
	"github.com/juggleim/jugglechat-server/storages"
)

func QryGroups(ctx context.Context, appkey string, offset string, limit int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Groups) {
	var startId int64 = 0
	var err error
	if offset != "" {
		startId, err = tools.DecodeInt(offset)
		if err != nil {
			startId = 0
		}
	}
	ret := &apimodels.Groups{
		Items: []*apimodels.Group{},
	}
	storage := storages.NewGroupStorage()
	grps, err := storage.QryGroups(appkey, startId, limit, isPositive)
	if err == nil {
		for _, grp := range grps {
			ret.Offset, _ = tools.EncodeInt(grp.ID)
			ret.Items = append(ret.Items, &apimodels.Group{
				GroupId:       grp.GroupId,
				GroupName:     grp.GroupName,
				GroupPortrait: grp.GroupPortrait,
				Owner: &apimodels.User{
					UserId: grp.CreatorId,
				},
				CreatedTime: grp.CreatedTime.UnixMilli(),
			})
		}
	}
	return errs.AdminErrorCode_Success, ret
}
