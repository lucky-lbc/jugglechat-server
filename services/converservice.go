package services

import (
	"context"
	"fmt"

	"github.com/juggleim/commons/ctxs"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/storages"
	"github.com/juggleim/jugglechat-server/storages/models"
)

type ConverConfItemKey string

const (
	ConverConfItemKey_MsgLifeTime ConverConfItemKey = "msg_life_time" //ms
)

func SetConverConfItem(ctx context.Context, targetId, subChannel string, converType int32, confItems map[string]interface{}) errs.IMErrorCode {
	appkey := ctxs.GetAppKeyFromCtx(ctx)
	userId := ctxs.GetRequesterIdFromCtx(ctx)
	fmt.Println("userId:", userId)
	fmt.Println("targetId:", targetId)
	converId := tools.GetConversationId(userId, targetId, converType)
	confs := []models.ConverConf{}
	for key, value := range confItems {
		if key == string(ConverConfItemKey_MsgLifeTime) {
			valueStr := fmt.Sprintf("%d", value)
			confs = append(confs, models.ConverConf{
				ConverId:   converId,
				ConverType: converType,
				SubChannel: subChannel,
				ItemKey:    key,
				ItemValue:  valueStr,
				ItemType:   0,
				AppKey:     appkey,
			})
		}
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
