package main

import (
	"fmt"

	"github.com/juggleim/commons/configures"
	"github.com/juggleim/commons/dbcommons"
	"github.com/juggleim/jugglechat-server/log"
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
}
