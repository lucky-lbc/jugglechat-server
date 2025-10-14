package dbmigrations

import (
	"bufio"
	"embed"
	"fmt"
	"sort"
	"strings"

	"github.com/lucky-lbc/commons/dbcommons"
	utils "github.com/lucky-lbc/commons/tools"
)

//go:embed sqls/*
var sqlFs embed.FS

const (
	JChatDbVersionKey       = "jchatdb_version"
	initVersion       int64 = 20250201
)

type CountResult struct {
	Count int64 `gorm:"count"`
}

func Upgrade() {
	// upgrade commons db
	dbcommons.Upgrade()
	// upgrade jchat db
	var currVersion int64 = 0
	dao := dbcommons.GlobalConfDao{}
	conf, err := dao.FindByKey(JChatDbVersionKey)
	if err == nil {
		ver, err := utils.String2Int64(conf.ConfValue)
		if err == nil && ver > 0 {
			currVersion = ver
		}
	} else {
		err = dao.Create(dbcommons.GlobalConfDao{
			ConfKey:   JChatDbVersionKey,
			ConfValue: fmt.Sprintf("%d", initVersion),
		})
		if err == nil {
			currVersion = initVersion
		}
	}
	fmt.Println("[JChatDbMigration]current version:", currVersion)
	sqlFiles, err := sqlFs.ReadDir("sqls")
	if err == nil {
		neededVers := []int64{}
		for _, sqlFile := range sqlFiles {
			fileName := sqlFile.Name()
			if len(fileName) == 12 {
				fileName = fileName[:8]
			}
			ver, err := utils.String2Int64(fileName)
			if err == nil && ver > 0 {
				neededVers = append(neededVers, ver)
			}
		}
		//sort
		sort.Slice(neededVers, func(i, j int) bool {
			return neededVers[i] < neededVers[j]
		})
		for _, ver := range neededVers {
			if ver > currVersion {
				sqlFileName := fmt.Sprintf("sqls/%d.sql", ver)
				fmt.Println("[DbMigration]start to execute sql file:", sqlFileName)
				err := executeSqlFile(sqlFileName)
				if err == nil {
					fmt.Println("[DbMigration]execute sql file success:", sqlFileName)
					dao.Upsert(dbcommons.GlobalConfDao{
						ConfKey:   JChatDbVersionKey,
						ConfValue: fmt.Sprintf("%d", ver),
					})
				}
			}
		}
	}
}

func executeSqlFile(fileName string) error {
	sqlFile, err := sqlFs.Open(fileName)
	if err != nil {
		fmt.Println("[DbMigration_Err]Read sql file err:", err, "file_name:", fileName)
		return err
	}
	defer sqlFile.Close()

	scanner := bufio.NewScanner(sqlFile)
	var queryBuilder strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") {
			continue
		}
		queryBuilder.WriteString(line)
		queryBuilder.WriteByte(' ')
		if strings.HasSuffix(line, ";") {
			query := strings.TrimSpace(queryBuilder.String())
			if query != "" {
				if err := dbcommons.GetDb().Exec(query).Error; err != nil {
					fmt.Println("[DbMigration_Err]Execute sql error:", err, query)
				}
			}
			queryBuilder.Reset()
		}
	}
	return nil
}
