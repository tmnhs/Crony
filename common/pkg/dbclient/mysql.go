package dbclient

import (
	"database/sql"
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _defaultDB *gorm.DB

func Init(dsn, logMode string, maxIdleConns, maxOpenConns int) (*gorm.DB, error) {

	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), setConfig(logMode)); err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetMaxOpenConns(maxOpenConns)
		_defaultDB = db
		return db, nil
	}
}

func GetMysqlDB() *gorm.DB {
	if _defaultDB == nil {
		logger.GetLogger().Error("mysql database is not initialized")
		return nil
	}
	return _defaultDB
}

func CreateDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}
