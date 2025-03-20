package services

import (
	"context"
	"fmt"

	"github.com/juggleim/jugglechat-server/apimodels"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services/imsdk"
	"github.com/juggleim/jugglechat-server/storages"
	"github.com/juggleim/jugglechat-server/storages/dbs"
	"github.com/juggleim/jugglechat-server/storages/models"
	"github.com/juggleim/jugglechat-server/utils"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func QryUserInfo(ctx context.Context, userId string) (errs.IMErrorCode, *apimodels.UserObj) {
	requestId := GetRequesterIdFromCtx(ctx)
	ret := &apimodels.UserObj{
		UserId: userId,
	}
	user := GetUser(ctx, userId)
	if user != nil {
		ret.Nickname = user.Nickname
		ret.Avatar = user.Avatar
	}
	if userId == requestId {
		ret.Settings = GetUserSettings(ctx, userId)
	} else {
		ret.IsFriend = checkFriend(ctx, requestId, userId)
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUserSettings(ctx context.Context, userId string) *apimodels.UserSettings {
	settings := &apimodels.UserSettings{}
	appkey := GetAppKeyFromCtx(ctx)
	storage := storages.NewUserExtStorage()
	exts, err := storage.QryExtFields(appkey, userId)
	if err == nil {
		for _, ext := range exts {
			if ext.ItemKey == apimodels.UserExtKey_Language {
				settings.Language = ext.ItemValue
			} else if ext.ItemKey == apimodels.UserExtKey_Undisturb {
				settings.Undisturb = ext.ItemValue
			} else if ext.ItemKey == apimodels.UserExtKey_FriendVerifyType {
				verifyType := utils.ToInt(ext.ItemValue)
				settings.FriendVerifyType = verifyType
			} else if ext.ItemKey == apimodels.UserExtKey_GrpVerifyType {
				verifyType := utils.ToInt(ext.ItemValue)
				settings.GrpVerifyType = verifyType
			}
		}
	}
	return settings
}

func SearchByPhone(ctx context.Context, phone string) (errs.IMErrorCode, *apimodels.Users) {
	requestId := GetRequesterIdFromCtx(ctx)
	appkey := GetAppKeyFromCtx(ctx)
	targetUserId := utils.ShortMd5(phone)
	storage := storages.NewUserStorage()
	user, err := storage.FindByPhone(appkey, phone)
	if err == nil && user != nil {
		targetUserId = user.UserId
	}
	users := &apimodels.Users{
		Items: []*apimodels.UserObj{},
	}
	users.Items = append(users.Items, &apimodels.UserObj{
		UserId:   user.UserId,
		Nickname: user.Nickname,
		Avatar:   user.UserPortrait,
		IsFriend: checkFriend(ctx, requestId, targetUserId),
	})
	return errs.IMErrorCode_SUCCESS, users
}

func UpdateUser(ctx context.Context, req *apimodels.UserObj) errs.IMErrorCode {
	appkey := GetAppKeyFromCtx(ctx)
	storage := storages.NewUserStorage()
	storage.Update(appkey, req.UserId, req.Nickname, req.Avatar)
	// sync to imserver
	sdk := imsdk.GetImSdk(appkey)
	if sdk != nil {
		sdk.Register(juggleimsdk.User{
			UserId:       req.UserId,
			Nickname:     req.Nickname,
			UserPortrait: req.Avatar,
		})
		if req.Nickname != "" {
			//update assistant
			sdk.AddBot(juggleimsdk.BotInfo{
				BotId:    GetAssistantId(req.UserId),
				Nickname: GetAssistantNickname(req.Nickname),
				Portrait: req.Avatar,
				BotType:  utils.IntPtr(int(apimodels.BotType_Custom)),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func UpdateUserSettings(ctx context.Context, req *apimodels.UserSettings) errs.IMErrorCode {
	appkey := GetAppKeyFromCtx(ctx)
	requestId := GetRequesterIdFromCtx(ctx)
	storage := storages.NewUserExtStorage()
	settings := map[juggleimsdk.UserSettingKey]string{}
	if req.Language != "" {
		storage.Upsert(models.UserExt{
			UserId:    requestId,
			ItemKey:   apimodels.UserExtKey_Language,
			ItemValue: req.Language,
			ItemType:  apimodels.AttItemType_Setting,
			AppKey:    appkey,
		})
		settings[juggleimsdk.UserSettingKey_Language] = req.Language
	}
	if req.Undisturb != "" {
		storage.Upsert(models.UserExt{
			UserId:    requestId,
			ItemKey:   apimodels.UserExtKey_Undisturb,
			ItemValue: req.Undisturb,
			ItemType:  apimodels.AttItemType_Setting,
			AppKey:    appkey,
		})
		settings[juggleimsdk.UserSettingKey_Undisturb] = req.Undisturb
	}
	storage.Upsert(models.UserExt{
		UserId:    requestId,
		ItemKey:   apimodels.UserExtKey_FriendVerifyType,
		ItemValue: utils.Int2String(int64(req.FriendVerifyType)),
		ItemType:  apimodels.AttItemType_Setting,
		AppKey:    appkey,
	})
	storage.Upsert(models.UserExt{
		UserId:    requestId,
		ItemKey:   apimodels.UserExtKey_GrpVerifyType,
		ItemValue: utils.Int2String(int64(req.GrpVerifyType)),
		ItemType:  apimodels.AttItemType_Setting,
		AppKey:    appkey,
	})
	//sync to im
	if len(settings) > 0 {
		if sdk := imsdk.GetImSdk(appkey); sdk != nil {
			sdk.SetUserSettings(juggleimsdk.User{
				UserId:   requestId,
				Settings: settings,
			})
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QueryMyGroups(ctx context.Context, limit int64, offset string) (errs.IMErrorCode, *apimodels.Groups) {
	appkey := GetAppKeyFromCtx(ctx)
	memberId := GetRequesterIdFromCtx(ctx)
	dao := dbs.GroupMemberDao{}
	var startId int64
	if offset != "" {
		startId, _ = utils.DecodeInt(offset)
	}
	groups, err := dao.QueryGroupsByMemberId(appkey, memberId, startId, limit)
	if err != nil {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	ret := &apimodels.Groups{
		Items: []*apimodels.Group{},
	}
	for _, group := range groups {
		ret.Offset, _ = utils.EncodeInt(group.ID)
		dao := dbs.GroupDao{}
		grpInfo, err := dao.FindById(appkey, group.GroupId)
		if err == nil {
			ret.Items = append(ret.Items, &apimodels.Group{
				GroupId:       grpInfo.GroupId,
				GroupName:     grpInfo.GroupName,
				GroupPortrait: grpInfo.GroupPortrait,
				// MemberCount:   grpInfo.MemberCount,
			})
		} else {
			fmt.Println(err)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func GetUser(ctx context.Context, userId string) *apimodels.UserObj {
	appkey := GetAppKeyFromCtx(ctx)
	u := &apimodels.UserObj{
		UserId: userId,
	}
	storage := storages.NewUserStorage()
	user, err := storage.FindByUserId(appkey, userId)
	if err == nil && user != nil {
		u.Nickname = user.Nickname
		u.Avatar = user.UserPortrait
		u.UserType = user.UserType
	}
	return u
}
