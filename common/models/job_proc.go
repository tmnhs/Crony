package models

import (
	"sync"
	"time"
)

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
	ID       int    `json:"id"` // pid
	JobID    int    `json:"jobId"`
	GroupId  int    `json:"group"`
	NodeUUID string `json:"node_uuid"`
	// parse from value
	JobProcVal

	Runnig int32
	HasPut int32
	Wg     sync.WaitGroup
}
