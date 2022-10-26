package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"gorm.io/gorm"
	"os"
)

func RegisterTables(db *gorm.DB) error {
	err := db.AutoMigrate(
		models.User{},
		models.Node{},
		models.Job{},
		models.JobLog{},
		models.Script{},
	)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("register table failed, error:%s", err.Error()))
		os.Exit(0)
	}
	entities := []models.User{
		{UserName: "root", Password: "e10adc3949ba59abbe56e057f20f883e", Role: models.RoleAdmin, Email: "333333333@qq.com"},
	}
	if exist := checkDataExist(db); !exist {
		if err := db.Table(models.CronyUserTableName).Create(&entities).Error; err != nil {
			return errors.Wrap(err, "Failed to initialize table data")
		}
	}
	logger.GetLogger().Info("register table success")
	return nil
}

func checkDataExist(db *gorm.DB) bool {
	if errors.Is(db.Table(models.CronyUserTableName).Where("username = ?", "root").First(&models.User{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
