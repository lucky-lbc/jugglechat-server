package services

import (
	"bytes"
	"context"
	"fmt"
	"jugglechat-server/apimodels"
	"jugglechat-server/errs"
	"jugglechat-server/services/aiengines"
	"jugglechat-server/services/imsdk"
	"jugglechat-server/storages"
	"jugglechat-server/utils"
	"time"

	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func AutoAnswer(ctx context.Context, req *apimodels.AssistantAnswerReq) (errs.IMErrorCode, *apimodels.AssistantAnswerResp) {
	if req == nil {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	if req.ChannelType == 1 {
		targetId := req.ConverId
		targetUser := GetUser(ctx, targetId)
		if targetUser.UserType == 1 {
			return errs.IMErrorCode_APP_DEFAULT, nil
		}
	}
	userId := GetRequesterIdFromCtx(ctx)
	promptStr := "你是一个智能回复生成器，能够根据用户提供的聊天记录，生成精彩回复。\n生成回复的一些限制条件：\n1. 只根据提供的聊天记录和上下文，生成回复，不进行无关的话题拓展；\n2. 确保回复的语音恰当、得体，不要产生冒犯性的表达；\n3. 回答简洁，不做过多延伸；\n4. 不要给我建议，直接以我的身份生成我该回复的内容；\n"
	appkey := GetAppKeyFromCtx(ctx)
	if req.PromptId != "" {
		pId, err := utils.DecodeInt(req.PromptId)
		if err == nil && pId > 0 {
			storage := storages.NewPromptStorage()
			prompt, err := storage.FindPrompt(appkey, userId, pId)
			if err == nil && prompt != nil && prompt.Prompts != "" {
				promptStr = promptStr + "回复的其他要求：" + prompt.Prompts
			}
		}
	}
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("下面是对话内容，格式是 {userid}:{聊天内容}， 请根据这个对话内容，以我的身份生成一条回复。\n")
	if len(req.Msgs) > 0 {
		for _, msg := range req.Msgs {
			if msg.SenderId != userId {
				buf.WriteString(fmt.Sprintf("%s:%s\n", msg.SenderId, msg.Content))
			} else {
				buf.WriteString(fmt.Sprintf("我:%s\n", msg.Content))
			}
		}
	} else {
		if req.ConverId == "" || req.ChannelType == 0 {
			return errs.IMErrorCode_APP_DEFAULT, nil
		}
		//qry history msg
		sdk := imsdk.GetImSdk(appkey)
		if sdk != nil {
			msgs, code, _, err := sdk.QryHisMsgs(userId, req.ConverId, juggleimsdk.ChannelType(req.ChannelType), 0, 5, false) //Chat-TODO  	MsgTypes:    []string{"jg:text"},
			if err == nil && code == 0 && msgs != nil {
				for _, msg := range msgs.Msgs {
					if msg.MsgType == "jg:text" {
						txtContent := &TextMsg{}
						err = utils.JsonUnMarshal([]byte(msg.MsgContent), txtContent)
						if err == nil {
							if msg.SenderId != userId {
								buf.WriteString(fmt.Sprintf("对方:%s\n", txtContent.Content))
							} else {
								buf.WriteString(fmt.Sprintf("我:%s\n", txtContent.Content))
							}
						}
					}
				}
			}
		}
	}
	if buf.Len() <= 0 {
		return errs.IMErrorCode_APP_DEFAULT, nil
	}
	content := buf.String()
	fmt.Println("----------------------------------------------------")
	fmt.Println(promptStr)
	fmt.Println("=")
	fmt.Println(content)
	fmt.Println("----------------------------------------------------")

	answer, streamMsgId := GenerateAnswer(ctx, promptStr, content, true)
	return errs.IMErrorCode_SUCCESS, &apimodels.AssistantAnswerResp{
		Answer:      answer,
		StreamMsgId: streamMsgId,
	}
}

func GenerateAnswer(ctx context.Context, prompt, question string, isSync bool) (string, string) {
	appkey := GetAppKeyFromCtx(ctx)
	userId := GetRequesterIdFromCtx(ctx)
	streamMsgId := utils.GenerateMsgId(time.Now().UnixMilli(), 0, "assistant")
	assistantInfo := aiengines.GetAiEngineInfo(ctx, appkey)
	if assistantInfo != nil && assistantInfo.AiEngine != nil {
		if isSync {
			buf := bytes.NewBuffer([]byte{})
			assistantInfo.AiEngine.StreamChat(ctx, userId, "assistant", prompt, question, func(answerPart string, isEnd bool) {
				if !isEnd {
					buf.WriteString(answerPart)
				}
			})
			return buf.String(), streamMsgId
		} else {
			go func() {
				assistantInfo.AiEngine.StreamChat(ctx, userId, "assistant", prompt, question, func(answerPart string, isEnd bool) {
					if !isEnd {
						sdk := imsdk.GetImSdk(appkey)
						if sdk != nil {
							sdk.SendSystemMsg(juggleimsdk.Message{
								SenderId:  "assistant",
								TargetIds: []string{userId},
								MsgType:   "jgs:aianswer",
								MsgContent: utils.ToJson(&StreamMsg{
									Content:     answerPart,
									StreamMsgId: streamMsgId,
								}),
								IsState: utils.BoolPtr(true),
							})
						}
					}
				})
			}()

			return "", streamMsgId
		}
	}
	return "No Answer", streamMsgId
}

type StreamMsg struct {
	Content     string `json:"content,omitempty"`
	StreamMsgId string `json:"stream_msg_id"`
}

type TextMsg struct {
	Content string `json:"content"`
}
