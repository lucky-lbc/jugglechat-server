package services

import (
	"context"
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/storages"
	"jugglechat-server/storages/models"
	"jugglechat-server/utils"
)

func QryAiBots(ctx context.Context, limit int64, offset string) (errs.IMErrorCode, *apimodels.AiBotInfos) {
	appkey := GetAppKeyFromCtx(ctx)
	storage := storages.NewBotConfStorage()
	var startId int64 = 0
	if offset != "" {
		intVal, err := utils.DecodeInt(offset)
		if err == nil {
			startId = intVal
		}
	}
	ret := &apimodels.AiBotInfos{
		Items: []*apimodels.AiBotInfo{},
	}
	items, err := storage.QryBotConfsWithStatus(appkey, models.BotStatus_Enable, startId, limit)
	if err == nil {
		for _, item := range items {
			ret.Offset, _ = utils.EncodeInt(item.ID)
			ret.Items = append(ret.Items, &apimodels.AiBotInfo{
				BotId:    item.BotId,
				Nickname: item.Nickname,
				Avatar:   item.BotPortrait,
				BotType:  int32(item.BotType),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
