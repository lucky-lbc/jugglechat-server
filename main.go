package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/commons/configures"
	"github.com/juggleim/commons/dbcommons"
	adminRouters "github.com/juggleim/jugglechat-server/admins/routers"
	"github.com/juggleim/jugglechat-server/log"
	"github.com/juggleim/jugglechat-server/routers"
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
	dbcommons.Upgrade()

	httpServer := gin.Default()
	routers.Route(httpServer, "jim")
	go httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))

	//start admin
	adminServer := gin.Default()
	adminRouters.Route(adminServer, "jconsole")
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
