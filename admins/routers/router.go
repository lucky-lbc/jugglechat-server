package routers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/commons/configures"
	"github.com/juggleim/commons/ctxs"
	"github.com/juggleim/jugglechat-server/admins/apis"
)

var Prefix string = ""

func RouteLogin(eng *gin.Engine, prefix string) *gin.RouterGroup {
	eng.Use(CorsHandler(), InjectCtx())
	group := eng.Group("/" + prefix)
	group.Use(apis.Validate)

	group.POST("/login", apis.Login)
	group.POST("/accounts/updpass", apis.UpdPassword)
	group.POST("/accounts/add", apis.AddAccount)
	group.POST("/accounts/delete", apis.DeleteAccounts)
	group.POST("/accounts/disable", apis.DisableAccounts)
	group.GET("/accounts/list", apis.QryAccounts)
	return group
}

func RouteProxy(group *gin.RouterGroup) *gin.RouterGroup {
	imAdminProxy := getImAdminProxy()
	if imAdminProxy != nil {
		group.GET("/apps/list", func(ctx *gin.Context) {
			imAdminProxy.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	return group
}

func Route(group *gin.RouterGroup) *gin.RouterGroup {
	//users
	group.GET("/apps/users/list", apis.QryUsers)
	group.POST("/apps/users/add")
	group.POST("/apps/users/update")
	group.POST("/apps/users/ban", apis.BanUsers)
	group.POST("/apps/users/unban", apis.UnBanUsers)

	//groups
	group.GET("/apps/groups/list", apis.QryGroups)

	//applications
	group.POST("/apps/applications/add", apis.AddApplication)
	group.POST("/apps/applications/update", apis.UpdApplication)
	group.POST("/apps/applications/delete", apis.DelApplications)
	group.GET("/apps/applications/list", apis.QryApplications)

	return group
}

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "*")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

func InjectCtx() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		appKey := ctx.Request.Header.Get("appkey")
		ctx.Set(string(ctxs.CtxKey_AppKey), appKey)
		ctx.Next()
	}
}

func getImAdminProxy() *httputil.ReverseProxy {
	if configures.Config.ImAdminDomain != "" {
		adminUrl, err := url.Parse(configures.Config.ImAdminDomain)
		if err == nil {
			proxy := httputil.NewSingleHostReverseProxy(adminUrl)
			proxy.Director = func(r *http.Request) {
				r.URL.Scheme = adminUrl.Scheme
				r.URL.Host = adminUrl.Host
				r.Host = adminUrl.Host

				r.Header.Set("X-Forwared-For", r.RemoteAddr)

				//rewrite path
				originalPath := r.URL.Path
				originalRawPath := r.URL.RawPath
				newPath, newRawPath := rewritePath(originalPath, originalRawPath)
				r.URL.Path = newPath
				if newRawPath != "" {
					r.URL.RawPath = newRawPath
				}
			}
			proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				fmt.Println("xxxx:", err)
				http.Error(w, "Internal error", http.StatusServiceUnavailable)
			}
			return proxy
		}
	}
	return nil
}

func rewritePath(path, rawPath string) (string, string) {
	newPath := strings.TrimPrefix(path, "/"+Prefix)
	if newPath == "" {
		newPath = "/"
	}
	newPath = "/admingateway" + newPath

	var newRawPath string
	if rawPath != "" {
		newRawPath = strings.TrimPrefix(rawPath, "/"+Prefix)
		if newRawPath == "" {
			newRawPath = "/"
		}
		newRawPath = "/admingateway" + newRawPath
	}
	return newPath, newRawPath
}
