package service

import (
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type ScriptService struct {
}

var DefaultScriptService = new(ScriptService)

func (script *ScriptService) Search(s *request.ReqScriptSearch) ([]models.Script, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyScriptTableName)
	if s.ID > 0 {
		db = db.Where("id = ?", s.ID)
	}
	if len(s.Name) > 0 {
		db = db.Where("name like ?", s.Name+"%")
	}
	scripts := make([]models.Script, 2)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&scripts).Error
	if err != nil {
		return nil, 0, err
	}
	return scripts, total, nil
}
