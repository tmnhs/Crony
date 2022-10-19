package service

import (
	"github.com/tmnhs/crony/common/models"
	"gorm.io/gorm"
)

func RegisterTables(db *gorm.DB) {
	_ = db.AutoMigrate(
		models.User{},
		models.Node{},
		models.Job{},
		models.JobLog{},
	)
}
