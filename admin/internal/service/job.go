package service

import (
	"github.com/tmnhs/crony/admin/internal/model"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type JobService struct {
}

var DefaultJobService = new(JobService)

func (j *JobService) Search(s *model.ReqJobSearch) ([]models.Job, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyJobTableName)
	if len(s.Name) > 0 {
		db = db.Where("name like ?", s.Name+"%")
	}
	if s.GroupId > 0 {
		db.Where("group_id = ?", s.GroupId)
	}
	if s.Type > 0 {
		db.Where("type = ?", s.Type)
	}
	if s.Kind > 0 {
		db.Where("kind = ?", s.Kind)
	}
	if s.Status > 0 {
		db.Where("status = ? ", s.Status)
	}
	jobs := make([]models.Job, 2)
	var total int64
	err := db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&jobs).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	return jobs, total, nil
}
