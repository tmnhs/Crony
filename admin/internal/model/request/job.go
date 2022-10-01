package request

import "github.com/tmnhs/crony/common/models"

type (
	ReqJobSearch struct {
		PageInfo
		Name    string         `json:"name" form:"name"`
		GroupId int            `json:"group_id" form:"group_id"`
		Kind    int            `json:"kind" form:"kind"`
		Type    models.JobType `json:"job_type" form:"type"`
		Status  int            `json:"status" form:"status"`
	}
	ReqJobLogSearch struct {
		PageInfo
		Name     string `json:"name" gorm:"name"`
		GroupId  int    `json:"group_id" gorm:"group_id"`
		JobId    int    `json:"job_id" gorm:"job_id"`
		NodeUUID string `json:"uuid" gorm:"node_uuid"`
		Success  *bool  `json:"success" gorm:"success"`
	}
)
