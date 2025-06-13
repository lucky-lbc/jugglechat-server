package apis

import (
	"net/http"
	"strconv"

	"github.com/juggleim/jugglechat-server/apis/models"
	"github.com/juggleim/jugglechat-server/apis/responses"
	"github.com/juggleim/jugglechat-server/ctxs"
	"github.com/juggleim/jugglechat-server/errs"
	"github.com/juggleim/jugglechat-server/services"
	"github.com/juggleim/jugglechat-server/services/aiengines"
	"github.com/juggleim/jugglechat-server/utils"

	"github.com/gin-gonic/gin"
)

func QryBots(ctx *gin.Context) {
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
	code, bots := services.QryAiBots(ctxs.ToCtx(ctx), int64(count), offset)
	if code != errs.IMErrorCode_SUCCESS {
		responses.ErrorHttpResp(ctx, code)
		return
	}
	responses.SuccessHttpResp(ctx, bots)
}

func BotMsgListener(ctx *gin.Context) {
	req := models.BotMsg{}
	if err := ctx.BindJSON(&req); err != nil {
		responses.ErrorHttpResp(ctx, errs.IMErrorCode_APP_REQ_BODY_ILLEGAL)
		return
	}
	prompt := `你是我的个人助理，帮我解答一切问题。`
	if req.Stream {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		assistantInfo := aiengines.GetAiEngineInfo(ctxs.ToCtx(ctx), ctx.GetString(string(ctxs.CtxKey_AppKey)))
		if assistantInfo != nil && assistantInfo.AiEngine != nil {
			idIndex := 1
			assistantInfo.AiEngine.StreamChat(ctxs.ToCtx(ctx), req.SenderId, req.BotId, prompt, req.Messages[0].Content, func(answerPart string, isEnd bool) {
				if !isEnd {
					item := &models.BotResponsePartData{
						Id:      utils.Int2String(int64(idIndex)),
						Type:    "message",
						Content: answerPart,
					}
					idIndex++
					ctx.Writer.Write([]byte("data: " + utils.ToJson(item) + "\n"))
					ctx.Writer.(http.Flusher).Flush()
				} else {
					ctx.Writer.Write([]byte("[DONE]"))
					ctx.Writer.(http.Flusher).Flush()
				}
			})
		}
	} else {

	}
}
