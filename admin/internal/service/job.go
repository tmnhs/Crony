package service

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
)

type JobService struct {
}

var DefaultJobService = new(JobService)

func (j *JobService) Search(s *request.ReqJobSearch) ([]models.Job, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyJobTableName)
	if s.ID > 0 {
		db = db.Where("id = ?", s.ID)
	}
	if len(s.Name) > 0 {
		db = db.Where("name like ?", s.Name+"%")
	}
	if len(s.RunOn) > 0 {
		db.Where("run_on = ?", s.RunOn)
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
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&jobs).Error
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
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Order("start_time desc").Find(&jobLogs).Error
	if err != nil {
		return nil, 0, err
	}

	return jobLogs, total, nil
}

//获取任务执行总数 1表示成功 0表示失败
func (j *JobService) GetTodayJobExcCount(success int) (int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyJobLogTableName).Where("start_time > ?", utils.GetTodayUnix()).Where("success = ?", success)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

//
func (j *JobService) GetRunningJobCount() (int64, error) {
	resp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdProcProfile), clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}

	return resp.Count, nil
}

const MaxJobCount = 10000

//优先分配工作任务最少的结点
func (j *JobService) AutoAllocateNode() string {
	//获取所有活着的节点
	nodeList := DefaultNodeWatcher.List2Array()

	resultCount, resultNodeUUID := MaxJobCount, ""
	for _, nodeUUID := range nodeList {
		count, err := DefaultNodeWatcher.GetJobCount(nodeUUID)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("node[%s] get job conut error:%s", nodeUUID, err.Error()))
			continue
		}
		if resultCount > count {
			resultCount, resultNodeUUID = count, nodeUUID
		}
	}
	return resultNodeUUID
}

//立即执行
func (j *JobService) Once(once *request.ReqJobOnce) (err error) {
	_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdOnce, once.JobId), once.NodeUUID)
	return
}
