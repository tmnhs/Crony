package request

import (
	"encoding/json"
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
	// default automatic assignment
	if r.Allocation == 0 {
		r.Allocation = models.AutoAllocation
	}
	notifyTo, _ := json.Marshal(r.NotifyToArray)
	r.NotifyTo = notifyTo
	scriptID, _ := json.Marshal(r.ScriptIDArray)
	r.ScriptID = scriptID
	return r.Check()
}
