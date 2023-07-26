package database

import (
	"time"

	"github.com/hoang-cao-long/learn-gorm/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGORM(config config.Config) (*gorm.DB, error) {
	mysqlConfig := config.DB.GetMysqlConfig()
	dsn := mysqlConfig.FormatDSN()
	dialect := mysql.Open(dsn)
	db, err := gorm.Open(dialect)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
