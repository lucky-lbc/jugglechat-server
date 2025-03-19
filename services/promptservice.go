package services

import (
	"context"
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/storages"
	"jugglechat-server/storages/models"
	"jugglechat-server/utils"

	"math"
)

func PromptAdd(ctx context.Context, req *apimodels.Prompt) (errs.IMErrorCode, *apimodels.Prompt) {
	storage := storages.NewPromptStorage()
	id, err := storage.Create(models.Prompt{
		UserId:  GetRequesterIdFromCtx(ctx),
		Prompts: req.Prompts,
		AppKey:  GetAppKeyFromCtx(ctx),
	})
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR, nil
	}
	idStr, _ := utils.EncodeInt(id)
	return errs.IMErrorCode_SUCCESS, &apimodels.Prompt{
		Id: idStr,
	}
}

func PromptUpdate(ctx context.Context, req *apimodels.Prompt) errs.IMErrorCode {
	id, _ := utils.DecodeInt(req.Id)
	if id <= 0 {
		return errs.IMErrorCode_APP_REQ_BODY_ILLEGAL
	}
	storage := storages.NewPromptStorage()
	err := storage.UpdatePrompts(GetAppKeyFromCtx(ctx), GetRequesterIdFromCtx(ctx), id, req.Prompts)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func PromptDel(ctx context.Context, req *apimodels.Prompt) errs.IMErrorCode {
	id, _ := utils.DecodeInt(req.Id)
	if id <= 0 {
		return errs.IMErrorCode_APP_REQ_BODY_ILLEGAL
	}
	storage := storages.NewPromptStorage()
	err := storage.DelPrompts(GetAppKeyFromCtx(ctx), GetRequesterIdFromCtx(ctx), id)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func PromptBatchDel(ctx context.Context, req *apimodels.PromptIds) errs.IMErrorCode {
	ids := []int64{}
	for _, idStr := range req.Ids {
		id, _ := utils.DecodeInt(idStr)
		if id > 0 {
			ids = append(ids, id)
		}
	}
	storage := storages.NewPromptStorage()
	err := storage.BatchDelPrompts(GetAppKeyFromCtx(ctx), GetRequesterIdFromCtx(ctx), ids)
	if err != nil {
		return errs.IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR
	}
	return errs.IMErrorCode_SUCCESS
}

func QryPrompts(ctx context.Context, count int64, offset string) (errs.IMErrorCode, *apimodels.Prompts) {
	var startId int64 = math.MaxInt64
	if offset != "" {
		id, _ := utils.DecodeInt(offset)
		if id > 0 {
			startId = id
		}
	}
	ret := &apimodels.Prompts{
		Items: []*apimodels.Prompt{},
	}
	storage := storages.NewPromptStorage()
	items, err := storage.QryPrompts(GetAppKeyFromCtx(ctx), GetRequesterIdFromCtx(ctx), count, startId)
	if err == nil {
		for _, item := range items {
			idStr, _ := utils.EncodeInt(item.ID)
			ret.Items = append(ret.Items, &apimodels.Prompt{
				Id:          idStr,
				Prompts:     item.Prompts,
				CreatedTime: item.CreatedTime,
			})
			ret.Offset = idStr
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
