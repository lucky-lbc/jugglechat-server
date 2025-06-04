package apis

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strconv"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
)

func CreateGroup(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := req.MemberIds
	if len(memberIds) <= 0 && len(req.GrpMembers) > 0 {
		ids := []string{}
		for _, member := range req.GrpMembers {
			ids = append(ids, member.UserId)
		}
		memberIds = ids
	}
	code, grpInfo := services.CreateGroup(services.ToCtx(ctx), &apimodels.GroupMembersReq{
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		MemberIds:     memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, &apimodels.Group{
		GroupId:       grpInfo.GroupId,
		GroupName:     grpInfo.GroupName,
		GroupPortrait: grpInfo.GroupPortrait,
	})
}

func UpdateGroup(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.UpdateGroup(services.ToCtx(ctx), &apimodels.GroupInfo{
		GroupId:       req.GroupId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func DissolveGroup(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DissolveGroup(services.ToCtx(ctx), req.GroupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QuitGroup(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.QuitGroup(services.ToCtx(ctx), req.GroupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func AddGrpMembers(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	memberIds := []string{}
	for _, user := range req.GrpMembers {
		memberIds = append(memberIds, user.UserId)
	}
	code := services.AddGrpMembers(services.ToCtx(ctx), &apimodels.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: memberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func DelGrpMembers(ctx *gin.Context) {
	req := apimodels.Group{}
	if err := ctx.BindJSON(&req); err != nil || req.GroupId == "" || len(req.MemberIds) <= 0 {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.DelGrpMembers(services.ToCtx(ctx), &apimodels.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QryGroupInfo(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	rpcCtx := services.ToCtx(ctx)
	code, grpInfo := services.QryGroupInfo(rpcCtx, groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, grpInfo)
}

func QryMyGroups(ctx *gin.Context) {
	offset := ctx.Query("offset")
	count := 20
	countStr := ctx.Query("count")
	var err error
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			count = 20
		}
	}
	code, grps := services.QueryMyGroups(services.ToCtx(ctx), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ret := &apimodels.Groups{
		Offset: grps.Offset,
		Items:  []*apimodels.Group{},
	}
	for _, grp := range grps.Items {
		ret.Items = append(ret.Items, &apimodels.Group{
			GroupId:       grp.GroupId,
			GroupName:     grp.GroupName,
			GroupPortrait: grp.GroupPortrait,
		})
	}
	SuccessHttpResp(ctx, ret)
}

func QryGrpMembers(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	offset := ctx.Query("offset")
	limit := 100
	limitStr := ctx.Query("limit")
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			limit = 100
		}
	}
	code, members := services.QueryGrpMembers(services.ToCtx(ctx), groupId, int64(limit), offset)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	ret := &apimodels.GroupMembersResp{
		Items: []*apimodels.GroupMember{},
	}
	for _, member := range members.Items {
		ret.Items = append(ret.Items, &apimodels.GroupMember{
			UserObj: apimodels.UserObj{
				UserId:   member.UserId,
				Nickname: member.Nickname,
				Avatar:   member.Avatar,
			},
		})
	}
	ret.Offset = members.Offset
	SuccessHttpResp(ctx, ret)
}

func CheckGroupMembers(ctx *gin.Context) {
	req := apimodels.CheckGroupMembersReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.CheckGroupMembers(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func SearchGroupMembers(ctx *gin.Context) {
	req := apimodels.SearchGroupMembersReq{}
	if err := ctx.BindJSON(&req); err != nil || req.GroupId == "" || req.Key == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code, resp := services.SearchGroupMembers(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, resp)
}

func SetGrpAnnouncement(ctx *gin.Context) {
	req := apimodels.GroupAnnouncement{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGrpAnnouncement(services.ToCtx(ctx), &apimodels.GroupAnnouncement{
		GroupId: req.GroupId,
		Content: req.Content,
	})
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func GetGrpAnnouncement(ctx *gin.Context) {
	groupId := ctx.Query("group_id")
	code, grpAnnounce := services.GetGrpAnnouncement(services.ToCtx(ctx), groupId)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, &apimodels.GroupAnnouncement{
		GroupId: grpAnnounce.GroupId,
		Content: grpAnnounce.Content,
	})
}

func SetGrpDisplayName(ctx *gin.Context) {
	req := apimodels.SetGroupDisplayNameReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	code := services.SetGrpDisplayName(services.ToCtx(ctx), &req)
	if code != errs.IMErrorCode_SUCCESS {
		ErrorHttpResp(ctx, code)
		return
	}
	SuccessHttpResp(ctx, nil)
}

func QryGrpQrCode(ctx *gin.Context) {
	grpId := ctx.Query("group_id")
	if grpId == "" {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	userId := ctx.GetString(string(services.CtxKey_RequesterId))

	m := map[string]interface{}{
		"action":   "join_group",
		"group_id": grpId,
		"user_id":  userId,
	}
	buf := bytes.NewBuffer([]byte{})
	qrCode, _ := qr.Encode(utils.ToJson(m), qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 400, 400)
	err := png.Encode(buf, qrCode)
	if err != nil {
		ErrorHttpResp(ctx, errs.IMErrorCode_APP_DEFAULT)
		return
	}
	SuccessHttpResp(ctx, map[string]string{
		"qr_code": base64.StdEncoding.EncodeToString(buf.Bytes()),
	})
}
