package aiengines

import (
	"context"
	"encoding/json"
	"jugglechat-server/storages"
	"jugglechat-server/storages/models"
	"sync"
)

type AssistantEngineType int

var (
	AssistantEngineType_SiliconFlow AssistantEngineType = 1
)

type IAiEngine interface {
	StreamChat(ctx context.Context, senderId, converId string, prompt string, question string, f func(answerPart string, isEnd bool))
}

var aiEngineCache *sync.Map
var aiEngineLock *sync.RWMutex

func init() {
	aiEngineCache = &sync.Map{}
	aiEngineLock = &sync.RWMutex{}
}

type AiEngineInfo struct {
	AppKey   string
	AiEngine IAiEngine
}

func GetAiEngineInfo(ctx context.Context, appkey string) *AiEngineInfo {
	key := appkey
	if val, exist := aiEngineCache.Load(key); exist {
		return val.(*AiEngineInfo)
	} else {
		aiEngineLock.Lock()
		defer aiEngineLock.Unlock()
		if val, exist := aiEngineCache.Load(key); exist {
			return val.(*AiEngineInfo)
		} else {
			aiEngineInfo := &AiEngineInfo{
				AppKey: appkey,
			}
			storage := storages.NewAiEngineStorage()
			ass, err := storage.FindEnableAiEngine(appkey)
			if err == nil {
				switch ass.EngineType {
				case models.EngineType_SiliconFlow:
					sfBot := &SiliconFlowEngine{}
					err = json.Unmarshal([]byte(ass.EngineConf), sfBot)
					if err == nil && sfBot.ApiKey != "" && sfBot.Url != "" && sfBot.Model != "" {
						aiEngineInfo.AiEngine = sfBot
					}
				}
			}
			aiEngineCache.Store(key, aiEngineInfo)
			return aiEngineInfo
		}
	}
}
