package services

import (
	"context"
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/storages"
	"jugglechat-server/storages/models"
	"jugglechat-server/utils"
	"time"
)

func PostAdd(ctx context.Context, req *apimodels.Post) (errs.IMErrorCode, *apimodels.Post) {
	appkey := GetAppKeyFromCtx(ctx)
	userId := GetRequesterIdFromCtx(ctx)
	storage := storages.NewPostStorage()
	postId := utils.GenerateUUIDString()
	createdTime := time.Now().UnixMilli()
	storage.Create(models.Post{
		PostId:      postId,
		Content:     []byte(utils.ToJson(req.Content)),
		UserId:      userId,
		CreatedTime: createdTime,
		UpdatedTime: time.Now(),
		AppKey:      appkey,
	})
	return errs.IMErrorCode_SUCCESS, &apimodels.Post{
		PostId:      postId,
		CreatedTime: createdTime,
	}
}

func QryPosts(ctx context.Context, startTime int64, limit int64, isPositive bool) (errs.IMErrorCode, *apimodels.Posts) {
	appkey := GetAppKeyFromCtx(ctx)
	ret := &apimodels.Posts{
		Items:      []*apimodels.Post{},
		IsFinished: true,
	}
	storage := storages.NewPostStorage()
	posts, err := storage.QryPosts(appkey, startTime, limit+1, isPositive)
	if err == nil {
		for _, p := range posts {
			postContent := &apimodels.PostContent{}
			utils.JsonUnMarshal(p.Content, postContent)
			post := &apimodels.Post{
				PostId:      p.PostId,
				Content:     postContent,
				UserInfo:    GetUser(ctx, p.UserId),
				CreatedTime: p.CreatedTime,
				UpdatedTime: p.UpdatedTime.UnixMilli(),
			}
			ret.Items = append(ret.Items, post)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
