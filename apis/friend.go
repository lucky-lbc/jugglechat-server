package apis

import (
	"strconv"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/gin-gonic/gin"
)

func QryFriends(ctx *gin.Context) {
	offset := ctx.Query("offset")
	count := 20
	var err error
	countStr := ctx.Query("count")
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, friends := services.QryFriends(services.ToCtx(ctx), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ret := &apimodels.Friends{
		Items:  []*apimodels.UserObj{},
		Offset: friends.Offset,
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &apimodels.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.Avatar,
			Pinyin:   friend.Pinyin,
			IsFriend: true,
		})
	}
	SuccessHttpResp(ctx, ret)
}

func QryFriendsWithPage(ctx *gin.Context) {
	var err error
	page := 1
	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}

	size := 20
	sizeStr := ctx.Query("size")
	if sizeStr != "" {
		size, err = strconv.Atoi(sizeStr)
		if err != nil {
			size = 20
		}
	}
	orderTag := ctx.Query("order_tag")
	code, friends := services.QryFriendsWithPage(services.ToCtx(ctx), int64(page), int64(size), orderTag)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ret := &apimodels.Friends{
		Items: []*apimodels.UserObj{},
	}
	for _, friend := range friends.Items {
		ret.Items = append(ret.Items, &apimodels.UserObj{
			UserId:   friend.UserId,
			Nickname: friend.Nickname,
			Avatar:   friend.Avatar,
			Pinyin:   friend.Pinyin,
			IsFriend: true,
		})
	}
	SuccessHttpResp(ctx, ret)
}

func SearchFriends(ctx *gin.Context) {
	req := apimodels.SearchFriendsReq{}
	if err := ctx.BindJSON(&req); err != nil || req.Key == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.SearchFriends(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func AddFriend(ctx *gin.Context) {
	req := apimodels.Friend{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.AddFriends(services.ToCtx(ctx), &apimodels.FriendIdsReq{
		FriendIds: []string{req.FriendId},
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func DelFriend(ctx *gin.Context) {
	req := apimodels.FriendIds{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelFriends(services.ToCtx(ctx), &apimodels.FriendIdsReq{
		FriendIds: req.FriendIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func ApplyFriend(ctx *gin.Context) {
	req := apimodels.ApplyFriend{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ApplyFriend(services.ToCtx(ctx), &apimodels.ApplyFriend{
		FriendId: req.FriendId,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func ConfirmFriend(ctx *gin.Context) {
	req := apimodels.ConfirmFriend{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.ConfirmFriend(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func MyFriendApplications(ctx *gin.Context) {
	startTimeStr := ctx.Query("start")
	start, err := utils.String2Int64(startTimeStr)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := utils.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyFriendApplications(services.ToCtx(ctx), start, int32(count), int32(order))
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func MyPendingFriendApplications(ctx *gin.Context) {
	startTimeStr := ctx.Query("start")
	start, err := utils.String2Int64(startTimeStr)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := utils.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryMyPendingFriendApplications(services.ToCtx(ctx), start, int32(count), int32(order))
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func FriendApplications(ctx *gin.Context) {
	startTimeStr := ctx.Query("start")
	start, err := utils.String2Int64(startTimeStr)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	countStr := ctx.Query("count")
	count, err := utils.String2Int64(countStr)
	if err != nil {
		count = 20
	} else {
		if count <= 0 || count > 50 {
			count = 20
		}
	}
	orderStr := ctx.Query("order")
	order, err := utils.String2Int64(orderStr)
	if err != nil || order > 1 || order < 0 {
		order = 0
	}
	code, resp := services.QryFriendApplications(services.ToCtx(ctx), start, count, int32(order))
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}
