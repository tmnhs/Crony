package service

import (
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type JobService struct {
}

var DefaultJobService = new(JobService)

func (j *JobService) Search(s *request.ReqJobSearch) ([]models.Job, int64, error) {
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
	if err != nil {
		return nil, 0, err
	}
	return jobs, total, nil
}

func (j *JobService) SearchJobLog(s *request.ReqJobLogSearch) ([]models.JobLog, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyJobLogTableName)
	if len(s.Name) > 0 {
		db = db.Where("name like ?", s.Name+"%")
	}
	if s.GroupId > 0 {
		db.Where("group_id = ?", s.GroupId)
	}
	if s.JobId > 0 {
		db.Where("job_id = ?", s.JobId)
	}
	if len(s.NodeUUID) > 0 {
		db.Where("node_uuid = ?", s.NodeUUID)
	}
	if s.Success != nil {
		db.Where("success = ? ", *s.Success)
	}
	jobLogs := make([]models.JobLog, 2)
	var total int64
	err := db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&jobLogs).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return jobLogs, total, nil
}
