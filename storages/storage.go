package storages

import (
	"jugglechat-server/storages/dbs"
	"jugglechat-server/storages/models"
)

func NewAppInfoStorage() models.IAppInfoStorage {
	return &dbs.AppInfoDao{}
}

func NewAppExtStorage() models.IAppExtStorage {
	return &dbs.AppExtDao{}
}

func NewUserStorage() models.IUserStorage {
	return &dbs.UserDao{}
}

func NewUserExtStorage() models.IUserExtStorage {
	return &dbs.UserExtDao{}
}

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}

func NewFriendApplicationStorage() models.IFriendApplicationStorage {
	return &dbs.FriendApplicationDao{}
}

func NewGrpApplicationStorage() models.IGrpApplicationStorage {
	return &dbs.GrpApplicationDao{}
}

func NewQrCodeRecordStorage() models.IQrCodeRecordStorage {
	return &dbs.QrCodeRecordDao{}
}

func NewSmsRecordStorage() models.ISmsRecordStorage {
	return &dbs.SmsRecordDao{}
}

func NewTeleBotStorage() models.ITeleBotStorage {
	return &dbs.TeleBotDao{}
}

func NewPromptStorage() models.IPromptStorage {
	return &dbs.PromptDao{}
}

func NewBotConfStorage() models.IBotConfStorage {
	return &dbs.BotConfDao{}
}

func NewAiEngineStorage() models.IAiEngineStorage {
	return &dbs.AiEngineDao{}
}

func NewGroupStorage() models.IGroupStorage {
	return &dbs.GroupDao{}
}

func NewGroupExtStorage() models.IGroupExtStorage {
	return &dbs.GroupExtDao{}
}

func NewGroupMemberStorage() models.IGroupMemberStorage {
	return &dbs.GroupMemberDao{}
}

func NewGroupAdminStorage() models.IGroupAdminStorage {
	return &dbs.GroupAdminDao{}
}

func NewPostStorage() models.IPostStorage {
	return &dbs.PostDao{}
}

func NewPostCommentStorage() models.IPostCommentStorage {
	return &dbs.PostCommentDao{}
}
