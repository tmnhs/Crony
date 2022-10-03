package service

import (
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type GroupService struct {
}

var DefaultGroupService = new(GroupService)

func (n *GroupService) Search(s *request.ReqGroupSearch) ([]models.Group, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyGroupTableName)
	if len(s.Name) > 0 {
		db = db.Where("name like ?", s.Name)
	}
	/*if s.Type > 0 {
		db.Where("type = ?", s.Type)
	}*/
	if s.ID > 0 {
		db.Where("id = ?", s.ID)
	}
	groups := make([]models.Group, 2)
	var total int64
	err := db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}
