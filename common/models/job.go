package models

import (
	"encoding/json"
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/common/pkg/utils/errors"
	"strings"
)

type JobType int

const (
	JobTypeCmd  = JobType(1)
	JobTypeHttp = JobType(2)

	HTTPMethodGet  = 1
	HTTPMethodPost = 2

	JobExcSuccess = 1
	JobExcFail    = 0

	JobStatusNotAssigned = 0
	JobStatusAssigned    = 1

	ManualAllocation = 1
	AutoAllocation   = 2
)

// 注册到 /crony/job/<node_uuid>/<job_id>
type Job struct {
	ID      int    `json:"id" gorm:"column:id"`
	Name    string `json:"name" gorm:"column:name" binding:"required"`
	Command string `json:"command" gorm:"column:command" binding:"required"`
	CmdUser string `json:"user" gorm:"column:cmd_user"`
	Timeout int64  `json:"timeout" gorm:"column:timeout"` // 任务执行时间超时设置，大于 0 时有效
	// 执行任务失败重试次数
	// 默认为 0，不重试
	RetryTimes int `json:"retry_times" gorm:"column:retry_times"`
	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	RetryInterval int64 `json:"retry_interval" gorm:"column:retry_interval"`
	// 任务类型
	Type       JobType `json:"job_type" gorm:"column:type" binding:"required"`
	HttpMethod int     `json:"http_method" gorm:"column:http_method"`
	// 执行失败是否发送通知
	NotifyType int `json:"notify_type" gorm:"column:notify_type"`
	//是否分配节点
	Status int `json:"status" gorm:"column:status"`
	// 发送通知地址
	NotifyTo      []byte `json:"-" gorm:"column:notify_to"`
	NotifyToArray []int  `json:"notify_to" gorm:"-"`
	Spec          string `json:"spec" gorm:"column:spec"`
	Note          string `json:"note" gorm:"column:note"`
	Created       int64  `json:"created" gorm:"column:created"`
	Updated       int64  `json:"updated" gorm:"column:updated"`
	// 执行任务的结点，用于记录 job log
	RunOn    string `json:"run_on" gorm:"column:run_on"`
	Hostname string `json:"host_name" gorm:"-"`
	Ip       string `json:"ip" gorm:"-"`
	// 用于存储分隔后的任务
	Cmd []string `json:"cmd" gorm:"-"`
}

func (j *Job) InitNodeInfo(status int, nodeUUID, hostname, ip string) {
	j.Status, j.RunOn, j.Hostname, j.Ip = status, nodeUUID, hostname, ip
}

func (j *Job) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyJobTableName).Create(j).Error
	if err == nil {
		insertId = j.ID
	}
	return
}

// 更新
func (j *Job) Update() error {
	return dbclient.GetMysqlDB().Table(CronyJobTableName).Updates(j).Error
}

func (j *Job) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where id = ?", CronyJobTableName), j.ID).Error
}

func (j *Job) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyJobTableName).Where("id = ? ", j.ID).First(j).Error
}
func (j *Job) Check() error {
	j.Name = strings.TrimSpace(j.Name)
	if len(j.Name) == 0 {
		return errors.ErrEmptyJobName
	}
	j.CmdUser = strings.TrimSpace(j.CmdUser)
	if j.RetryTimes == 0 {
		j.RetryTimes = 3
	}
	if j.RetryInterval == 0 {
		j.RetryTimes = 1
	}
	if j.Timeout == 0 {
		j.RetryTimes = 5
	}
	if len(strings.TrimSpace(j.Command)) == 0 {
		return errors.ErrEmptyJobCommand
	}
	if len(j.Cmd) == 0 && j.Type == JobTypeCmd {
		j.SplitCmd()
	}
	return nil
}

func (j *Job) SplitCmd() {
	ps := strings.SplitN(j.Command, " ", 2)
	if len(ps) == 1 {
		j.Cmd = ps
		return
	}
	j.Cmd = make([]string, 0, 2)
	j.Cmd = append(j.Cmd, ps[0])
	j.Cmd = append(j.Cmd, utils.ParseCmdArguments(ps[1])...)
}

func (j *Job) Val() string {
	data, err := json.Marshal(j)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
