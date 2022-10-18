package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type JobLog struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	JobId    int    `json:"job_id" gorm:"column:job_id"`
	Command  string `json:"command" gorm:"column:command"`
	IP       string `json:"ip" gorm:"column:ip"` // node ip
	Hostname string `json:"hostname" gorm:"column:hostname"`
	NodeUUID string `json:"node_uuid" gorm:"column:node_uuid"`
	Success  bool   `json:"success" gorm:"column:success"`

	Output string `json:"output" gorm:"column:output"`
	Spec   string `json:"spec" gorm:"column:spec"`

	RetryTimes int   `json:"retry_times" gorm:"column:retry_times"`
	StartTime  int64 `json:"start_time" gorm:"column:start_time"`
	EndTime    int64 `json:"end_time" gorm:"column:end_time"`
}

func (jb *JobLog) Update() error {
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
