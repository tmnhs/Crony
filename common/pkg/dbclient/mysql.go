package dbclient

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _defaultDB *gorm.DB

func Init(dsn, logMode string, maxIdleConns, maxOpenConns int) (*gorm.DB, error) {

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置
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

func Insert(table string, val interface{}) error {
	if _defaultDB == nil {
		return errors.ErrClientNotFound
	}
	return _defaultDB.Table(table).Create(val).Error
}

func DeleteById(table string, id int64) error {
	if _defaultDB == nil {
		return errors.ErrClientNotFound
	}
	return _defaultDB.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", table), id).Error
}
