package main

import (
	"fmt"

	"github.com/juggleim/commons/configures"
	"github.com/juggleim/commons/dbcommons"
	"github.com/juggleim/commons/imsdk"
	"github.com/juggleim/commons/tools"
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

	sdk := imsdk.GetImSdk("nsw3sue72begyv7y")

	resp, code, _, err := sdk.QryGlobalConvers(0, 10)
	fmt.Println(err)
	fmt.Println(code)
	fmt.Println(tools.ToJson(resp))
}
