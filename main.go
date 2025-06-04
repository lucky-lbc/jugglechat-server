package main

import (
	"fmt"
	"net/http"

	"github.com/juggleim/jugglechat-server/apis"
	"github.com/juggleim/jugglechat-server/configures"
	"github.com/juggleim/jugglechat-server/log"
	"github.com/juggleim/jugglechat-server/storages/dbs/dbcommons"

	"github.com/gin-gonic/gin"
)

func main() {
	//init configure
	if err := configures.InitConfigures(); err != nil {
		fmt.Println("Init Configures failed", err)
		return
	}
	//init log
	log.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		log.Error("Init Mysql failed.", err)
		return
	}

	server := gin.Default()
	server.Use(CorsHandler())
	group := server.Group("/jim")
	group.Use(apis.Validate)

	group.POST("/login", apis.Login)
	group.GET("/login/qrcode", apis.GenerateQrCode)
	group.POST("/login/qrcode/check", apis.CheckQrCode)
	group.POST("/sms/send", apis.SmsSend)
	group.POST("/sms_login", apis.SmsLogin)
	group.POST("/sms/login", apis.SmsLogin)
	group.POST("/email/send", apis.EmailSend)
	group.POST("/email/login", apis.EmailLogin)
	group.POST("/login/qrcode/confirm", apis.ConfirmQrCode)
	group.POST("/file_cred", apis.GetFileCred)

	group.GET("/bots/list", apis.QryBots)

	group.POST("/assistants/answer", apis.AssistantAnswer)
	group.POST("/assistants/prompts/add", apis.PromptAdd)
	group.POST("/assistants/prompts/update", apis.PromptUpdate)
	group.POST("/assistants/prompts/del", apis.PromptDel)
	group.POST("/assistants/prompts/batchdel", apis.PromptBatchDel)
	group.GET("/assistants/prompts/list", apis.QryPrompts)

	group.POST("/bots/messages/listener", apis.BotMsgListener)

	group.POST("/users/update", apis.UpdateUser)
	group.POST("/users/updsettings", apis.UpdateUserSettings)
	group.POST("/users/search", apis.SearchByPhone)
	group.GET("/users/info", apis.QryUserInfo)
	group.GET("/users/qrcode", apis.QryUserQrCode)

	group.POST("/telegrambots/add", apis.TelegramBotAdd)
	group.POST("/telegrambots/del", apis.TelegramBotDel)
	group.POST("/telegrambots/batchdel", apis.TelegramBotBatchDel)
	group.GET("/telegrambots/list", apis.TelegramBotList)

	group.POST("/groups/add", apis.CreateGroup)
	group.POST("/groups/create", apis.CreateGroup)
	group.POST("/groups/update", apis.UpdateGroup)
	group.POST("/groups/dissolve", apis.DissolveGroup)
	group.POST("/groups/members/add", apis.AddGrpMembers)
	group.POST("/groups/apply", apis.GroupApply)
	group.POST("/groups/invite", apis.GroupInvite)
	group.POST("/groups/quit", apis.QuitGroup)
	group.POST("/groups/members/del", apis.DelGrpMembers)
	group.GET("/groups/members/list", apis.QryGrpMembers)
	group.POST("/groups/members/check", apis.CheckGroupMembers)
	group.POST("/groups/members/search", apis.SearchGroupMembers)
	group.GET("/groups/info", apis.QryGroupInfo)
	group.GET("/groups/qrcode", apis.QryGrpQrCode)
	group.POST("/groups/setgrpannouncement", apis.SetGrpAnnouncement)
	group.GET("/groups/getgrpannouncement", apis.GetGrpAnnouncement)
	group.POST("/groups/setdisplayname", apis.SetGrpDisplayName)
	//group manage
	group.POST("/groups/management/chgowner", apis.ChgGroupOwner)
	group.POST("/groups/management/administrators/add", apis.AddGrpAdministrator)
	group.POST("/groups/management/adminstrators/del", apis.DelGrpAdministrator)
	group.GET("/groups/management/administrators/list", apis.QryGrpAdministrators)
	group.POST("/groups/management/setmute", apis.SetGroupMute)
	group.POST("/groups/management/setgrpverifytype", apis.SetGrpVerifyType)
	group.POST("/groups/management/sethismsgvisible", apis.SetGrpHisMsgVisible)
	group.GET("/groups/mygroups", apis.QryMyGroups)
	// grp application
	group.GET("/groups/myapplications", apis.QryMyGrpApplications)
	group.GET("/groups/mypendinginvitations", apis.QryMyPendingGrpInvitations)
	group.GET("/groups/grpinvitations", apis.QryGrpInvitations)
	group.GET("/groups/grppendingapplications", apis.QryGrpPendingApplications)

	group.GET("/friends/list", apis.QryFriendsWithPage)
	group.POST("/friends/search", apis.SearchFriends)
	group.POST("/friends/add", apis.AddFriend)
	group.POST("/friends/apply", apis.ApplyFriend)
	group.POST("/friends/confirm", apis.ConfirmFriend)
	group.POST("/friends/del", apis.DelFriend)
	group.GET("/friends/applications", apis.FriendApplications)
	group.GET("/friends/myapplications", apis.MyFriendApplications)
	group.GET("/friends/mypendingapplications", apis.MyPendingFriendApplications)

	//post
	group.GET("/posts/list", apis.QryPosts)
	group.GET("/posts/info", apis.PostInfo)
	group.POST("/posts/add", apis.PostAdd)
	group.POST("/posts/update")
	group.POST("/posts/del")
	group.POST("/posts/reactions/add")
	group.GET("/posts/reactions/list")

	group.GET("/postcomments/list", apis.QryPostComments)
	group.POST("/postcomments/add", apis.PostCommentAdd)
	group.POST("/postcomments/update")
	group.POST("/postcomments/del")

	fmt.Println("Start Server with port:", configures.Config.Port)
	server.Run(fmt.Sprintf(":%d", configures.Config.Port))

}
func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
