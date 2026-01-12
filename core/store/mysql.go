// @Title
// @Description
// @Author  Wangwengang  2023/12/12 23:58
// @Update  Wangwengang  2023/12/12 23:58
package store

import (
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"github.com/wwengg/threego/core/store/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormMysqlByConfig 初始化Mysql数据库用过传入配置
func GormMysqlByConfig(m sconfig.Mysql, log slog.Slog) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m.Prefix, m.Singular, log)); err != nil {
		panic(err)
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
