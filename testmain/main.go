package main

import (
	"fmt"
	"github.com/lucky-lbc/jugglechat-server/storages/dbs"

	"github.com/lucky-lbc/jugglechat-server/commons/configures"
	"github.com/lucky-lbc/jugglechat-server/commons/dbcommons"
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

	dao := dbs.UserDao{}
	users, err := dao.QryUsers("appkey", "ser", 0, 10, false)
	fmt.Println(err, users)
}
