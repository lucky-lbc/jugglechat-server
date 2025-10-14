package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lucky-lbc/commons/configures"
	"github.com/lucky-lbc/commons/dbcommons"
	adminRouters "github.com/lucky-lbc/jugglechat-server/admins/routers"
	"github.com/lucky-lbc/jugglechat-server/log"
	"github.com/lucky-lbc/jugglechat-server/routers"
	"github.com/lucky-lbc/jugglechat-server/storages/dbs/dbmigrations"
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
	//upgrade db
	dbmigrations.Upgrade()

	httpServer := gin.Default()
	routers.Route(httpServer, "jim")
	go httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))

	//start admin
	adminServer := gin.Default()
	group := adminRouters.RouteLogin(adminServer, "admingateway")
	adminRouters.RouteProxy(group)
	adminRouters.Route(group)
	adminRouters.LoadJuggleChatAdminWeb(adminServer)
	go adminServer.Run(fmt.Sprintf(":%d", configures.Config.AdminPort))

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		signal.Stop(sigChan)
		close(closeChan)
	}()

	<-closeChan
}
