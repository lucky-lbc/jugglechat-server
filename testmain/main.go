package main

import (
	"fmt"
	"time"

	"github.com/juggleim/jugglechat-server/configures"
	"github.com/juggleim/jugglechat-server/log"
	"github.com/juggleim/jugglechat-server/storages/dbs"
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

	userDao := dbs.UserDao{}

	var startId int64 = 0
	for {
		users, err := userDao.QryUsers("nsw3sue72begyv7y", startId, 100)
		if err == nil {
			for _, user := range users {
				startId = user.ID
				userDao.Update(user.AppKey, user.UserId, user.Nickname, "")
				fmt.Println("xxxx:", user.UserId, user.Nickname)
				time.Sleep(500 * time.Millisecond)
			}
		}
		if len(users) < 100 {
			break
		}
	}
}
