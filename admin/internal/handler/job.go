package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"time"
)

type JobRouter struct {
}

var defaultJobRouter = new(JobRouter)

func (j *JobRouter) CreateOrUpdate(c *gin.Context) {
	var req request.ReqJobUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[create_job] request parameter error", c)
		return
	}
	if err := req.Valid(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("create_job check error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorJobFormat, "[create_job] check error", c)
		return
	}
	var err error
	var insertId int
	t := time.Now()
	if len(req.NotifyToArray) > 0 {
		notifyTo, _ := json.Marshal(req.NotifyToArray)
		req.NotifyTo = notifyTo
	}

	if req.Allocation == models.AutoAllocation {
		if !config.GetConfigModels().System.CmdAutoAllocation && req.Type == models.JobTypeCmd {
			resp.FailWithMessage(resp.ERROR, "[create_job] The shell command is not supported to automatically assign nodes by default.", c)
			return
		}
		//自动分配
		nodeUUID := service.DefaultJobService.AutoAllocateNode()
		if nodeUUID == "" {
			logger.GetLogger().Error(fmt.Sprintf("[create_job] auto allocate node error"))
			resp.FailWithMessage(resp.ERROR, "[create_job] auto allocate node error", c)
			return
		}
		req.RunOn = nodeUUID
	} else if req.Allocation == models.ManualAllocation {
		//手动分配
		if len(req.RunOn) == 0 {
			resp.FailWithMessage(resp.ERROR, "[create_job] manually assigned node can't be null", c)
			return
		}
		node := &models.Node{UUID: req.RunOn}
		_ = node.FindByUUID()
		if node.Status == models.NodeConnFail {
			resp.FailWithMessage(resp.ERROR, "[create_job] manually assigned node inactivation", c)
			return
		}
	}
	if req.ID > 0 {
		//update
		job := &models.Job{ID: req.ID}
		_ = job.FindById()
		oldNodeUUID := job.RunOn
		if oldNodeUUID != "" {
			_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, oldNodeUUID, req.ID))
			if err != nil {
				logger.GetLogger().Error(fmt.Sprintf("[update_job] delete etcd node[%s]  error:%s", oldNodeUUID, err.Error()))
				resp.FailWithMessage(resp.ERROR, "[update_job] delete etcd node error", c)
				return
			}
		}
		req.Updated = t.Unix()
		err = req.Update()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[update_job] into db  error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[update_job] into db id error", c)
			return
		}
	} else {
		//create
		req.Created = t.Unix()
		insertId, err = req.Insert()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[create_job] insert job into db error:%s", err.Error()))
			resp.FailWithMessage(resp.ERROR, "[create_job] insert job into db error", c)
			return
		}
		req.ID = insertId
	}
	b, err := json.Marshal(req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_job] json marshal job error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[create_job] json marshal job error", c)
		return
	}
	//添加至etcd
	_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, req.RunOn, req.ID), string(b))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_job] etcd put job error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[create_job] etcd put job error", c)
		return
	}

	resp.OkWithDetailed(req, "operate success", c)
}

func (j *JobRouter) Delete(c *gin.Context) {
	var req request.ByIDS
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_job] request parameter error", c)
		return
	}
	for _, id := range req.IDS {
		//先查找再删除etcd之后再删除数据库
		job := models.Job{ID: id}
		err := job.FindById()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[delete_job] find job by id :%d error:%s", id, err.Error()))
			continue
		}
		_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, job.RunOn, id))
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[delete_job] etcd delete job error:%s", err.Error()))
			continue
		}
		err = job.Delete()
		if err != nil {
			logger.GetLogger().Error(fmt.Sprintf("[delete_job] into db error:%s", err.Error()))
			continue
		}
	}
	resp.OkWithMessage("delete success", c)
}

func (j *JobRouter) FindById(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[find_job] request parameter error", c)
		return
	}

	//先查找再删除etcd之后再删除数据库
	job := models.Job{ID: req.ID}
	err := job.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[find_job] find job by id :%d error:%s", req.ID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[find_job] find job by id error", c)
		return
	}
	if len(job.NotifyTo) != 0 {
		_ = json.Unmarshal(job.NotifyTo, &job.NotifyToArray)
	}
	resp.OkWithDetailed(job, "find success", c)
}

func (j *JobRouter) Search(c *gin.Context) {
	var req request.ReqJobSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_job] request parameter error", c)
		return
	}
	req.Check()
	jobs, total, err := service.DefaultJobService.Search(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_job] search job error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[search_job] search job error", c)
		return
	}
	var resultJobs []models.Job
	for _, job := range jobs {
		_ = json.Unmarshal(job.NotifyTo, &job.NotifyToArray)
		resultJobs = append(resultJobs, job)
	}
	resp.OkWithDetailed(resp.PageResult{
		List:     resultJobs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "search success", c)
}

func (j *JobRouter) SearchLog(c *gin.Context) {
	var req request.ReqJobLogSearch
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_job_log] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[search_job_log] request parameter error", c)
		return
	}
	req.Check()
	jobs, total, err := service.DefaultJobService.SearchJobLog(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[search_job_log] db error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[search_job_log] db error", c)
		return
	}
	resp.OkWithDetailed(resp.PageResult{
		List:     jobs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "search success", c)
}

//手动执行
func (j *JobRouter) Once(c *gin.Context) {
	var req request.ReqJobOnce
	var err error
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[job_once] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[job_once] request parameter error", c)
		return
	}
	jobModel := &models.Job{ID: req.JobId}
	err = jobModel.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[job_once] job_id[%d] not exist db:%s", req.JobId, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[job_once] job not exist ", c)
		return
	}
	err = service.DefaultJobService.Once(&req)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[job_once] etcd put job_id :%d error:%s", req.JobId, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[job_once] etcd put  error", c)
		return
	}
	resp.OkWithMessage("job once success", c)
}
