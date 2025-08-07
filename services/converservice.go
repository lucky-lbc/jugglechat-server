package services

import (
	"context"

	"github.com/juggleim/commons/ctxs"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/tools"
	apimodels "github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/storages"
	"github.com/juggleim/jugglechat-server/storages/models"
)

type ConverConfItemKey string

const (
	ConverConfItemKey_MsgLifeTime ConverConfItemKey = "msg_life_time" //ms
)

func SetConverConfItem(ctx context.Context, targetId, subChannel string, converType int32, items []*apimodels.ConverConfItem) errs.IMErrorCode {
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	userId := ctxs.GetRequesterIdFromCtx(ctx)
	converId := tools.GetConversationId(userId, targetId, converType)
	confs := []models.ConverConf{}
	for _, item := range items {
		confs = append(confs, models.ConverConf{
			ConverId:   converId,
			ConverType: converType,
			SubChannel: subChannel,
			ItemKey:    item.ItemKey,
			ItemValue:  item.ItemValue,
			ItemType:   item.ItemType,
			AppKey:     appkey,
		})
	}
	storage := storages.NewConverConfStorage()
	err := storage.BatchUpsert(confs)
	if err != nil {
		return errs.IMErrorCode_APP_DEFAULT
	}
	return errs.IMErrorCode_SUCCESS
}

func GetConverConfItems(ctx context.Context, targetId, subChannel string, converType int32) (errs.IMErrorCode, map[string]interface{}) {
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	userId := ctxs.GetRequesterIdFromCtx(ctx)
	converId := tools.GetConversationId(userId, targetId, converType)
	ret := map[string]interface{}{}
	storage := storages.NewConverConfStorage()
	confMap, err := storage.QryConverConfs(appkey, converId, subChannel, converType)
	if err == nil {
		for _, conf := range confMap {
			if conf.ItemKey == string(ConverConfItemKey_MsgLifeTime) {
				var msgLifeTime int64 = 0
				i, err := tools.String2Int64(conf.ItemValue)
				if err == nil && i > 0 {
					msgLifeTime = i
				}
				ret[conf.ItemKey] = msgLifeTime
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
