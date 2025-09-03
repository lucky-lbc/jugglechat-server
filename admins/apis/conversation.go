package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/juggleim/commons/errs"
	"github.com/juggleim/commons/imsdk"
	"github.com/juggleim/commons/responses"
	"github.com/juggleim/commons/tools"
	"github.com/juggleim/jugglechat-server/admins/apis/models"
	"github.com/juggleim/jugglechat-server/admins/services"
)

func QryConversations(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		responses.AdminErrorHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		count, err := tools.String2Int64(startStr)
		if err == nil && count > 0 {
			start = count
		}
	}

	var count int64 = 20
	countStr := ctx.Query("count")
	if countStr != "" {
		countVal, err := tools.String2Int64(countStr)
		if err == nil && countVal > 0 {
			count = countVal
		}
	}
	ret := &models.GlobalConversations{
		Items: []*models.GlobalConversation{},
	}
	sdk := imsdk.GetImSdk(appkey)
	if sdk != nil {
		resp, code, _, err := sdk.QryGlobalConvers(start, int(count))
		if err == nil && code == 200 && resp != nil {
			for _, item := range resp.Items {
				conver := &models.GlobalConversation{
					ChannelType: item.ChannelType,
					Sender:      services.QryUserInfo(appkey, item.UserId),
					Time:        item.Time,
				}
				if item.ChannelType == 1 {
					conver.Receiver = services.QryUserInfo(appkey, item.TargetId)
				} else if item.ChannelType == 2 {
					conver.Group = services.QryGroupInfo(appkey, item.TargetId)
				}
				ret.Items = append(ret.Items, conver)
			}
		}
	}
	responses.AdminSuccessHttpResp(ctx, ret)
}
