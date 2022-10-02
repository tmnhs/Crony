package models

import (
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

	KindAlone = 1

	//job log  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT "1->成功 2->正在运行 3->失败",
	JobLogStatusSuccess = 1
	JobLogStatusProcess = 2
	JobLogStatusFail    = 3

	ManualAllocation = 1
	AutoAllocation   = 2
)

// 需要执行的 cron cmd 命令
// 注册到 /cronsun/cmd/<node_uuid>/<group_id>/<job_id>
type Job struct {
	ID      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name" binding:"required"`
	GroupId int    `json:"group_id" gorm:"-" `
	Command string `json:"command" gorm:"command" binding:"required"`
	CmdUser string `json:"user" gorm:"cmd_user"`
	Pause   bool   `json:"pause" gorm:"-"`         // 可手工控制的状态
	Timeout int64  `json:"timeout" gorm:"timeout"` // 任务执行时间超时设置，大于 0 时有效
	// 设置任务在单个节点上可以同时允许多少个
	// 针对两次任务执行间隔比任务执行时间要长的任务启用
	Parallels int64 `json:"parallels" gorm:"-"`
	// 执行任务失败重试次数
	// 默认为 0，不重试
	RetryTimes int `json:"retry_times" gorm:"retry_times"`
	// 执行任务失败重试时间间隔
	// 单位秒，如果不大于 0 则马上重试
	RetryInterval int64 `json:"retry_interval" gorm:"retry_interval"`
	// 任务类型
	// 0: 普通任务
	// 1: 单机任务
	// 如果为单机任务，node 加载任务的时候 Parallels 设置 1
	Kind       int     `json:"kind" gorm:"kind"`
	Type       JobType `json:"job_type" gorm:"type" binding:"required"`
	HttpMethod int     `json:"http_method" gorm:"http_method"`
	// 执行失败是否发送通知
	NotifyStatus bool `json:"notify_status" gorm:"notify_status"`
	NotifyType   int  `json:"notify_type" gorm:"notify_type"`
	Status       int  `json:"status" gorm:"status"`
	// 发送通知地址
	NotifyTo      []byte `json:"-" gorm:"notify_to"`
	NotifyToArray []int  `json:"notify_to" gorm:"-"`
	NotifyToType  int    `json:"notify_to_type" gorm:"notify_to_type"`
	Spec          string `json:"spec" gorm:"spec"`

	Created int64 `json:"created" gorm:"created"`
	Updated int64 `json:"updated" gorm:"updated"`
	// 平均执行时间，单位 ms
	AvgTime int64 `json:"avg_time" gorm:"-"`
	// 单独对任务指定日志清除时间
	LogExpiration int `json:"log_expiration" gorm:"-"`

	// 执行任务的结点，用于记录 job log
	RunOn    string `json:"run_on" gorm:"-"`
	Hostname string `json:"host_name" gorm:"-"`
	Ip       string `json:"ip" gorm:"-"`
	// 用于存储分隔后的任务
	Cmd []string `json:"cmd" gorm:"-"`
}

func (j *Job) InitNodeInfo(nodeUUID, hostname, ip string) {
	j.RunOn, j.Hostname, j.Ip = nodeUUID, hostname, ip
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
	if j.LogExpiration < 0 {
		j.LogExpiration = 0
	}

	j.CmdUser = strings.TrimSpace(j.CmdUser)

	// 不修改 Command 的内容，简单判断是否为空
	if len(strings.TrimSpace(j.Command)) == 0 {
		return errors.ErrEmptyJobCommand
	}
	if len(j.Cmd) == 0 {
		j.SplitCmd()
	}
	//todo 安全性

	//security := conf.Config.Security
	//if !security.Open {
	//	return nil
	//}
	//
	/*if len(conf.Config.Security.Users) == 0 {
		return true
	}

	for _, u := range conf.Config.Security.Users {
		if j.User == u {
			return true
		}
	}*/
	//
	//if !j.validCmd() {
	//	return ErrSecurityInvalidCmd
	//}

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
