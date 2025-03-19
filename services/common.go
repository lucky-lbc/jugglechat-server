package services

import (
	"context"

	"github.com/gin-gonic/gin"
)

type CtxKey string

const (
	CtxKey_AppKey      CtxKey = "CtxKey_AppKey"
	CtxKey_Session     CtxKey = "CtxKey_Session"
	CtxKey_RequesterId CtxKey = "CtxKey_RequesterId"
)

func ToCtx(ginCtx *gin.Context) context.Context {
	rpcCtx := context.Background()
	rpcCtx = context.WithValue(rpcCtx, CtxKey_AppKey, ginCtx.GetString(string(CtxKey_AppKey)))
	rpcCtx = context.WithValue(rpcCtx, CtxKey_Session, ginCtx.GetString(string(CtxKey_Session)))
	currentUserId := ginCtx.GetString(string(CtxKey_RequesterId))
	if currentUserId != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_RequesterId, currentUserId)
	}
	return rpcCtx
}

func GetAppKeyFromCtx(ctx context.Context) string {
	if appKey, ok := ctx.Value(CtxKey_AppKey).(string); ok {
		return appKey
	}
	return ""
}

func GetRequesterIdFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_RequesterId).(string); ok {
		return requesterId
	}
	return ""
}
