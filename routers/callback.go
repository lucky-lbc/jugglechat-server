package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/jugglechat-server/apis"
	"github.com/lucky-lbc/jugglechat-server/services"
)

func CallbackRoute(eng *gin.Engine, prefix string) *gin.RouterGroup {
	services.CallbackUrlPrefix = prefix
	group := eng.Group("/" + prefix)

	group.POST("/assistants/msgcallback", apis.AssistantMsgCallback)

	return group
}
