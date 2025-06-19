package ctxs

import (
	"context"

	"github.com/gin-gonic/gin"
)

type CtxKey string

const (
	CtxKey_AppKey      CtxKey = "CtxKey_AppKey"
	CtxKey_Session     CtxKey = "CtxKey_Session"
	CtxKey_RequesterId CtxKey = "CtxKey_RequesterId"

	CtxKey_Account CtxKey = "CtxKey_Account"
)

func ToCtx(ginCtx *gin.Context) context.Context {
	rpcCtx := context.Background()
	appkey := ginCtx.GetString(string(CtxKey_AppKey))
	if appkey != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_AppKey, appkey)
	}
	session := ginCtx.GetString(string(CtxKey_Session))
	if session != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_Session, session)
	}
	currentUserId := ginCtx.GetString(string(CtxKey_RequesterId))
	if currentUserId != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_RequesterId, currentUserId)
	}
	account := ginCtx.GetString(string(CtxKey_Account))
	if account != "" {
		rpcCtx = context.WithValue(rpcCtx, CtxKey_Account, account)
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

func GetSessionFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_Session).(string); ok {
		return requesterId
	}
	return ""
}

func GetAccountFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_Account).(string); ok {
		return requesterId
	}
	return ""
}
