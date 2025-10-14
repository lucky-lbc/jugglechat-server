package main

import (
	"fmt"

	"github.com/lucky-lbc/commons/configures"
	"github.com/lucky-lbc/commons/dbcommons"
	"github.com/lucky-lbc/jugglechat-server/log"
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
