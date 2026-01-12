// @Title
// @Description
// @Author  Wangwengang  2023/12/12 23:57
// @Update  Wangwengang  2023/12/12 23:57
package store

import (
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"gorm.io/gorm"
)

func DBList(config *[]sconfig.SpecializedDB, log slog.Slog) map[string]*gorm.DB {
	dbMap := make(map[string]*gorm.DB)
	for _, info := range *config {
		if info.Disable {
			continue
		}
		switch info.Type {
		case "mysql":
			if info.GeneralDB.LogZap {
				dbMap[info.AliasName] = GormMysqlByConfig(sconfig.Mysql{GeneralDB: info.GeneralDB}, log)
			} else {
				dbMap[info.AliasName] = GormMysqlByConfig(sconfig.Mysql{GeneralDB: info.GeneralDB}, nil)
			}

		//case "mssql":
		//	dbMap[info.AliasName] = GormMssqlByConfig(config.Mssql{GeneralDB: info.GeneralDB})
		//case "pgsql":
		//	dbMap[info.AliasName] = GormPgSqlByConfig(config.Pgsql{GeneralDB: info.GeneralDB})
		//case "oracle":
		//	dbMap[info.AliasName] = GormOracleByConfig(config.Oracle{GeneralDB: info.GeneralDB})
		default:
			continue
		}
	}
	return dbMap
}
