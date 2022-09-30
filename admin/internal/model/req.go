package model

import "github.com/tmnhs/crony/common/models"

type (
	// PageInfo Paging common input parameter structure
	PageInfo struct {
		Page     int `json:"page" form:"page"`           // 页码
		PageSize int `json:"page_size" form:"page_size"` // 每页大小
	}
	ByID struct {
		ID int `json:"id" form:"id"`
	}
	ReqUserLogin struct {
		UserName string `json:"username" form:"username" binding:"required,min=2,max=20"`
		Password string `json:"password" form:"password" binding:"required,min=4,max=20,alphanum"`
	}
	ReqUserRegister struct {
		UserName string `json:"username" form:"username" binding:"required,min=2,max=20"`
		Password string `json:"password" form:"password" binding:"required,min=4,max=20,alphanum"`
	}
	ReqJobSearch struct {
		PageInfo
		Name    string         `json:"name" form:"name"`
		GroupId int            `json:"group_id" form:"group_id"`
		Kind    int            `json:"kind" form:"kind"`
		Type    models.JobType `json:"job_type" form:"type"`
		Status  int            `json:"status" form:"status"`
	}
)

func (page *PageInfo) Check() {
	if page.PageSize <= 0 {
		page.PageSize = 20
	}
	if page.Page <= 0 {
		page.Page = 1
	}
}
