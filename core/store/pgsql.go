// @Title
// @Description
// @Author  Wangwengang  2023/12/12 23:58
// @Update  Wangwengang  2023/12/12 23:58
package store

import (
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"github.com/wwengg/threego/core/store/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GormPgSqlByConfig 初始化PostgreSQL数据库通过传入配置
func GormPgSqlByConfig(p sconfig.Pgsql, log slog.Slog) *gorm.DB {
	if p.Dbname == "" {
		return nil
	}
	if db, err := gorm.Open(postgres.Open(p.Dsn()), internal.Gorm.Config(p.Prefix, p.Singular, log)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		return db
	}
}
