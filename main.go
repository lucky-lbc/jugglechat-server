package main

import (
	"fmt"
	"jugglechat-server/apis"
	"jugglechat-server/configures"
	"jugglechat-server/log"
	"jugglechat-server/storages/dbs/dbcommons"
	"net/http"

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
	group.POST("/groups/info", apis.QryGroupInfo)
	group.POST("/groups/qrcode", apis.QryGrpQrCode)
	group.POST("/groups/setgrpannouncement", apis.SetGrpAnnouncement)
	group.GET("/groups/getgrpannouncement", apis.GetGrpAnnouncement)
	group.POST("/groups/setdisplayname", apis.SetGrpDisplayName)
	//group manage
	group.POST("/groups/management/chgowner", apis.ChgGroupOwner)
	group.POST("/groups/management/administrators/add", apis.AddGrpAdministrator)
	group.POST("/groups/management/adminstrators/del", apis.DelGrpAdministrator)
	group.GET("/groups/management/administrators/list", apis.QryGrpAdministrators)
	group.POST("/groups/management/setmute", apis.SetGroupMute)

	group.POST("/jim/groups/management/setmute", apis.SetGroupMute)
	group.POST("/jim/groups/management/setgrpverifytype", apis.SetGrpVerifyType)
	group.POST("/jim/groups/management/sethismsgvisible", apis.SetGrpHisMsgVisible)
	group.GET("/jim/groups/mygroups", apis.QryMyGroups)
	// grp application
	group.GET("/jim/groups/myapplications", apis.QryMyGrpApplications)
	group.GET("/jim/groups/mypendinginvitations", apis.QryMyPendingGrpInvitations)
	group.GET("/jim/groups/grpinvitations", apis.QryGrpInvitations)
	group.GET("/jim/groups/grppendingapplications", apis.QryGrpPendingApplications)

	group.GET("/jim/friends/list", apis.QryFriendsWithPage)
	group.POST("/jim/friends/add", apis.AddFriend)
	group.POST("/jim/friends/apply", apis.ApplyFriend)
	group.POST("/jim/friends/confirm", apis.ConfirmFriend)
	group.POST("/jim/friends/del", apis.DelFriend)
	group.GET("/jim/friends/applications", apis.FriendApplications)
	group.GET("/jim/friends/myapplications", apis.MyFriendApplications)
	group.GET("/jim/friends/mypendingapplications", apis.MyPendingFriendApplications)

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
