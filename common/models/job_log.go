package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

//日志输出
type JobLog struct {
	ID       int    `json:"id" gorm:"id"`
	Name     string `json:"name" gorm:"name"`
	GroupId  int    `json:"group_id" gorm:"group_id"`
	JobId    int    `json:"job_id" gorm:"job_id"`
	Command  string `json:"command" gorm:"command"`
	IP       string `json:"ip" gorm:"ip"` // node ip
	Hostname string `json:"hostname" gorm:"hostname"`
	NodeUUID string `json:"uuid" gorm:"node_uuid"`
	Success  bool   `json:"success" gorm:"success"`

	Output string `json:"output" gorm:"output"`
	Spec   string `json:"spec" gorm:"spec"`

	// 执行任务失败重试次数
	// 默认为 0，不重试
	RetryTimes int   `json:"retry_times" gorm:"retry_times"`
	StartTime  int64 `json:"start_time" gorm:"start_time"`
	EndTime    int64 `json:"end_time" gorm:"end_time"`
}

// 更新
func (jb *JobLog) Update() error {
	//只会更新非零字段
	return dbclient.GetMysqlDB().Table(CronyJobLogTableName).Updates(jb).Error
}

func (jb *JobLog) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where id = ?", CronyJobLogTableName), jb.ID).Error
}

func (jb *JobLog) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyJobLogTableName).Create(jb).Error
	if err == nil {
		insertId = jb.ID
	}
	return
}
