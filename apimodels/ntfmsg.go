package apimodels

var GroupNotifyMsgType string = "jgd:grpntf"

type GroupNotify struct {
	Operator *UserObj        `json:"operator"`
	Name     string          `json:"name"`
	Members  []*UserObj      `json:"members"`
	Type     GroupNotifyType `json:"type"`
}

type GroupNotifyType int

const (
	GroupNotifyType_AddMember    GroupNotifyType = 1
	GroupNotifyType_RemoveMember GroupNotifyType = 2
	GroupNotifyType_Rename       GroupNotifyType = 3
	GroupNotifyType_ChgOwner     GroupNotifyType = 4
	GroupNotifyType_Join         GroupNotifyType = 5
)

var FriendNotifyMsgType string = "jgd:friendntf"

type FriendNotify struct {
	Type int `json:"type"`
}

var SystemFriendApplyConverId string = "friend_apply"
var FriendApplicationMsgType string = "jgd:friendapply"

type FriendApplyNotify struct {
	SponsorId   string `json:"sponsor_id"`
	RecipientId string `json:"recipient_id"`
}
