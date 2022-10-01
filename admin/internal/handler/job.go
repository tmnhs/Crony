package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/admin/internal/model/resp"
	"github.com/tmnhs/crony/admin/internal/service"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"time"
)

type JobRouter struct {
}

var defaultJobRouter = new(JobRouter)

func (j *JobRouter) CreateOrUpdate(c *gin.Context) {
	var req models.Job
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[create_job] request parameter error", c)
		return
	}
	//todo node是否存活
	if err := req.Check(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("create_job check error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorJobFormat, "[create_job] error error", c)
		return
	}
	logger.GetLogger().Debug(fmt.Sprintf("create job req:%#v", req))
	var err error
	var insertId int
	t := time.Now()
	req.Created = time.Now().Unix()
	notifyTo, _ := json.Marshal(req.NotifyToArray)
	req.NotifyTo = notifyTo
	//想更改数据库
	if req.ID > 0 {
		//update
		insertId = req.ID
		req.Updated = t.UnixNano()
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
	//todo 分配方法：手动和自动
	_, err = etcdclient.Put(fmt.Sprintf(etcdclient.KeyEtcdJob, req.RunOn, req.GroupId, req.ID), string(b))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[create_job] etcd put job error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[create_job] etcd put job error", c)
		return
	}

	resp.OkWithMessage("add success", c)
}

func (j *JobRouter) Delete(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_job] request parameter error", c)
		return
	}
	//先查找再删除etcd之后再删除数据库
	job := models.Job{ID: req.ID}
	err := job.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] find job by id :%d error:%s", req.ID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_job] find job by id error", c)
		return
	}
	_, err = etcdclient.Delete(fmt.Sprintf(etcdclient.KeyEtcdJob, job.RunOn, job.GroupId, req.ID))
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] etcd delete job error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_job] etcd delete job error", c)
		return
	}
	err = job.Delete()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] into db error:%s", err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_job] into db error", c)
		return
	}
	resp.OkWithMessage("delete success", c)
}

func (j *JobRouter) FindById(c *gin.Context) {
	var req request.ByID
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] request parameter error:%s", err.Error()))
		resp.FailWithMessage(resp.ErrorRequestParameter, "[delete_job] request parameter error", c)
		return
	}
	//先查找再删除etcd之后再删除数据库
	//job, err := models.FindJobById(req.ID)
	job := models.Job{ID: req.ID}
	err := job.FindById()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("[delete_job] find job by id :%d error:%s", req.ID, err.Error()))
		resp.FailWithMessage(resp.ERROR, "[delete_job] find job by id error", c)
		return
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
	resp.OkWithDetailed(resp.PageResult{
		List:     jobs,
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
