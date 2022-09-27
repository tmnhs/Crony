package models

import (
	"sync"
	"time"
)

type JobType int

const (
	JobTypeCmd  = JobType(1)
	JobTypeHttp = JobType(2)

	HTTPMethodGet  = 1
	HTTPMethodPost = 2

	KindAlone = 1
)

// 需要执行的 cron cmd 命令
// 注册到 /cronsun/cmd/groupName/<id>
type Job struct {
	ID      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	GroupId int    `json:"group" gorm:"gid"`
	Command string `json:"command" gorm:"command"`
	CmdUser string `json:"user" gorm:"cmd_user"`
	Pause   bool   `json:"pause"`                  // 可手工控制的状态
	Timeout int    `json:"timeout" gorm:"timeout"` // 任务执行时间超时设置，大于 0 时有效
	// 设置任务在单个节点上可以同时允许多少个
	// 针对两次任务执行间隔比任务执行时间要长的任务启用
	Parallels int64 `json:"parallels"`
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
	Type       JobType `json:"job_type" gorm:"type"`
	HttpMethod int     `json:"http_method" gorm:"http_method"`
	// 执行失败是否发送通知
	NotifyStatus bool `json:"notify_status" gorm:"notify_status"`
	NotifyType   bool `json:"notify_type" gorm:"notify_type"`
	Status       int  `json:"status" gorm:"status"`
	// 发送通知地址
	NotifyTo     []byte `json:"notify_to" gorm:"-"`
	NotifyToType []byte `json:"notify_to_type" gorm:"-"`
	Spec         string `json:"spec" gorm:"spec"`

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
	// 控制同时执行任务数
	Count *int64 `json:"-"`
}

//日志输出
type JobLog struct {
	ID       int    `json:"id" gorm:"id"`
	Name     string `json:"name" gorm:"name"`
	GroupId  int    `json:"group" gorm:"gid"`
	JobId    int    `json:"jid" gorm:"jid"`
	Command  string `json:"command" gorm:"command"`
	IP       string `json:"ip" gorm:"ip"` // node ip
	Hostname string `json:"hostname" gorm:"hostname"`
	NodeUUID string `json:"uuid" gorm:"node_uuid"`
	Status   int    `json:"status" gorm:"status"`
	Output   int    `json:"output" gorm:"output"`
	Spec     string `json:"spec" gorm:"spec"`

	// 执行任务失败重试次数
	// 默认为 0，不重试
	RetryTimes int   `json:"retry_times" gorm:"retry_times"`
	StartTime  int64 `json:"start_time" gorm:"start_time"`
	EndTime    int64 `json:"end_time" gorm:"end_time"`
}

type jobLink struct {
	gname string
	// rule id
	rules map[string]bool
}
type Link map[string]map[string]*jobLink

type JobProcVal struct {
	Time   time.Time `json:"time"`   // 开始执行时间
	Killed bool      `json:"killed"` // 是否强制杀死
}

// 当前执行中的任务信息
// key: /cronsun/proc/node/group/jobId/pid
// value: 开始执行时间
// key 会自动过期，防止进程意外退出后没有清除相关 key，过期时间可配置
type JobProc struct {
	// parse from key path
	ID     string `json:"id"` // pid
	JobID  string `json:"jobId"`
	Group  string `json:"group"`
	NodeID string `json:"nodeId"`
	// parse from value
	JobProcVal

	Runnig int32
	HasPut int32
	Wg     sync.WaitGroup
}

func (j *Job) InitNodeInfo(nodeID, hostname, ip string) {
	var c int64
	j.Count, j.RunOn, j.Hostname, j.Ip = &c, nodeID, hostname, ip
}
