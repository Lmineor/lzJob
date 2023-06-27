package lzJob

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

func initMysql(cfg *Mysql) *gorm.DB {
	if cfg.Db == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       cfg.Dsn(), // DSN data source name
		DefaultStringSize:         191,       // string 类型字段的默认长度
		SkipInitializeWithVersion: false,     // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		klog.Fatalf("open mysql connection failed %s", err)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		return db
	}
}
