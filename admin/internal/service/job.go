package service

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"time"
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

//获取今日任务执行总数 1表示成功 0表示失败
func (j *JobService) GetTodayJobExcCount(success int) (int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyJobLogTableName).Where("start_time > ? and end_time!=0 and success = ?", utils.GetTodayUnix(), success)
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

//某个时间段内每天的任务数量
func (j *JobService) GetJobExcCount(start, end int64, success int) ([]resp.DateCount, error) {
	var dateCount []resp.DateCount
	db := dbclient.GetMysqlDB().Table(models.CronyJobLogTableName).Select("FROM_UNIXTIME( start_time, '%Y-%m-%d' ) AS date", "COUNT( * ) AS count ").Group("date").Order("date ASC").Where("start_time > ? and start_time<?  and end_time!=0 and success = ?", start, end, success)
	err := db.Find(&dateCount).Error
	if err != nil {
		return nil, err
	}
	return dateCount, nil
}

func (j *JobService) GetNotAssignedJob() (jobs []models.Job, err error) {
	jobs = make([]models.Job, 2)
	err = dbclient.GetMysqlDB().Table(models.CronyJobTableName).Where("status = ?", models.JobStatusNotAssigned).Find(&jobs).Error
	return
}

//
func (j *JobService) GetRunningJobCount() (int64, error) {
	wresp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdProcProfile), clientv3.WithPrefix(), clientv3.WithCountOnly())
	if err != nil {
		return 0, err
	}

	return wresp.Count, nil
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
	//默认存在时间60秒
	_, err = etcdclient.PutWithTtl(fmt.Sprintf(etcdclient.KeyEtcdOnce, once.JobId), once.NodeUUID, 60)
	return
}

func RunLogCleaner(cleanPeriod time.Duration, expiration int64) (close chan struct{}) {
	t := time.NewTicker(cleanPeriod)
	close = make(chan struct{})
	go func() {
		for {
			select {
			case <-t.C:
				err := cleanupLogs(expiration)
				if err != nil {
					logger.GetLogger().Error(fmt.Sprintf("clean up logs at time:%v error:%s", time.Now(), err.Error()))
				}
			case <-close:
				return
			}
		}
	}()
	return
}

func cleanupLogs(expirationTime int64) error {
	sql := fmt.Sprintf("delete from %s where start_time < ?", models.CronyJobLogTableName)
	return dbclient.GetMysqlDB().Exec(sql, time.Now().Unix()-expirationTime).Error
}
