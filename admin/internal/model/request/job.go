package request

import (
	"github.com/tmnhs/crony/common/models"
)

type (
	ReqJobSearch struct {
		PageInfo
		ID     int            `json:"id" form:"id"`
		Name   string         `json:"name" form:"name"`
		RunOn  string         `json:"run_on" form:"run_on"`
		Type   models.JobType `json:"job_type" form:"type"`
		Status int            `json:"status" form:"status"`
	}
	ReqJobLogSearch struct {
		PageInfo
		Name     string `json:"name" form:"name"`
		JobId    int    `json:"job_id" form:"job_id"`
		NodeUUID string `json:"node_uuid" form:"node_uuid"`
		Success  *bool  `json:"success" form:"success"`
	}
	ReqJobUpdate struct {
		*models.Job
		//分配方式
		Allocation int `json:"allocation" form:"allocation" binding:"required"`
	}
	ReqJobOnce struct {
		JobId    int    `json:"job_id" form:"job_id"`
		NodeUUID string `json:"node_uuid" form:"node_uuid"`
	}
	ReqJobKill struct {
		JobId    int    `json:"job_id" form:"job_id"`
		NodeUUID string `json:"node_uuid" form:"node_uuid"`
	}
)

func (r *ReqJobUpdate) Valid() error {
	//默认自动分配
	if r.Allocation == 0 {
		r.Allocation = models.AutoAllocation
	}
	return r.Check()
}
