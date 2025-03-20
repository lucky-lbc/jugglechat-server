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
			//top comments
			code, topComments := QryPostComments(ctx, p.PostId, 0, 10, false)
			if code == errs.IMErrorCode_SUCCESS && topComments != nil {
				post.TopComments = topComments.Items
			}
			ret.Items = append(ret.Items, post)
		}
		if len(ret.Items) > int(limit) {
			ret.Items = ret.Items[:limit]
			ret.IsFinished = false
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryPostInfo(ctx context.Context, postId string) (errs.IMErrorCode, *apimodels.Post) {
	appkey := GetAppKeyFromCtx(ctx)
	ret := &apimodels.Post{}
	storage := storages.NewPostStorage()
	post, err := storage.FindById(appkey, postId)
	if err == nil {
		postContent := &apimodels.PostContent{}
		err = utils.JsonUnMarshal(post.Content, postContent)
		if err == nil {
			ret.Content = postContent
		}
		ret.PostId = post.PostId
		ret.CreatedTime = post.CreatedTime
		ret.UpdatedTime = post.UpdatedTime.UnixMilli()
		ret.UserInfo = GetUser(ctx, post.UserId)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func PostCommentAdd(ctx context.Context, req *apimodels.PostComment) (errs.IMErrorCode, *apimodels.PostComment) {
	appkey := GetRequesterIdFromCtx(ctx)
	userId := GetRequesterIdFromCtx(ctx)
	postId := req.PostId
	storage := storages.NewPostCommentStorage()
	commentId := utils.GenerateUUIDString()
	currTime := time.Now()
	parentUserId := req.ParentUserId
	if parentUserId == "" && req.ParentUserInfo != nil && req.ParentUserInfo.UserId != "" {
		parentUserId = req.ParentUserInfo.UserId
	}
	storage.Create(models.PostComment{
		CommentId:       commentId,
		PostId:          postId,
		ParentCommentId: req.ParentCommentId,
		ParentUserId:    parentUserId,
		Text:            req.Text,
		UserId:          userId,
		CreatedTime:     currTime.UnixMilli(),
		UpdatedTime:     currTime,
		AppKey:          appkey,
	})
	return errs.IMErrorCode_SUCCESS, nil
}

func QryPostComments(ctx context.Context, postId string, startTime int64, limit int64, isPositive bool) (errs.IMErrorCode, *apimodels.PostComments) {
	appkey := GetAppKeyFromCtx(ctx)
	ret := &apimodels.PostComments{
		Items:      []*apimodels.PostComment{},
		IsFinished: true,
	}
	storage := storages.NewPostCommentStorage()
	comments, err := storage.QryPostComments(appkey, postId, startTime, limit+1, isPositive)
	if err == nil {
		for _, c := range comments {
			ret.Items = append(ret.Items, &apimodels.PostComment{
				PostId:          c.PostId,
				CommentId:       c.CommentId,
				Text:            c.Text,
				ParentCommentId: c.ParentCommentId,
				ParentUserInfo:  GetUser(ctx, c.ParentUserId),
				UserInfo:        GetUser(ctx, c.UserId),
				CreatedTime:     c.CreatedTime,
				UpdatedTime:     c.UpdatedTime.UnixMilli(),
			})
		}
		if len(ret.Items) > int(limit) {
			ret.Items = ret.Items[:limit]
			ret.IsFinished = false
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
