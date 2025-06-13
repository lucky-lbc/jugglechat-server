package services

import (
	"context"
	"fmt"

	apimodels "github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/configures"
	"github.com/juggleim/jugglechat-server/services/imsdk"
	"github.com/juggleim/jugglechat-server/utils"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func InitUserAssistant(ctx context.Context, userId, nickname, portrait string) {
	appkey := GetAppKeyFromCtx(ctx)
	sdk := imsdk.GetImSdk(appkey)
	if sdk != nil {
		if appinfo, exist := GetAppInfo(appkey); exist {
			apikey, err := GenerateApiKey(appkey, appinfo.AppSecureKey)
			if err == nil {
				botId := GetAssistantId(userId)
				callbackUrl := configures.Config.AiBotCallbackUrl
				if callbackUrl == "" {
					if configures.Config.ConnectManager.WsProxyPort != 0 {
						callbackUrl = fmt.Sprintf("http://127.0.0.1:%d/jim/bots/messages/listener", configures.Config.ConnectManager.WsProxyPort)
					} else if configures.Config.ConnectManager.WsPort != 0 {
						callbackUrl = fmt.Sprintf("http://127.0.0.1:%d/jim/bots/messages/listener", configures.Config.ConnectManager.WsPort)
					}
				}
				sdk.AddBot(juggleimsdk.BotInfo{
					BotId:    botId,
					Nickname: GetAssistantNickname(nickname),
					Portrait: portrait,
					BotType:  utils.IntPtr(int(apimodels.BotType_Custom)),
					BotConf:  fmt.Sprintf(`{"url":"%s","api_key":"%s","bot_id":"%s"}`, callbackUrl, apikey, botId),
				})
				sdk.SendPrivateMsg(juggleimsdk.Message{
					SenderId:       botId,
					TargetIds:      []string{userId},
					MsgType:        "jg:text",
					MsgContent:     `{"content":"欢迎注册，我是您的私人助理，任何问题都可以问我！"}`,
					IsNotifySender: utils.BoolPtr(false),
				})
			}
		}
	}
}

func GetAssistantId(userId string) string {
	botId := fmt.Sprintf("ass_%s", userId)
	return botId
}

func GetAssistantNickname(nickname string) string {
	return fmt.Sprintf("%s 的助理", nickname)
}
