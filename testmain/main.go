package main

import (
	"fmt"

	"github.com/juggleim/jugglechat-server/configures"
	"github.com/juggleim/jugglechat-server/log"
	"github.com/juggleim/jugglechat-server/storages/dbs/dbcommons"
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
