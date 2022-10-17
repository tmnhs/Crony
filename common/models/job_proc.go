package models

import (
	"encoding/json"
	"sync"
	"time"
)

type JobProcVal struct {
	Time   time.Time `json:"time"`   // 开始执行时间
	Killed bool      `json:"killed"` // 是否强制杀死
}

type JobProc struct {
	// parse from key path
	ID       int    `json:"id"` // pid
	JobID    int    `json:"job_id"`
	NodeUUID string `json:"node_uuid"`
	// parse from value
	JobProcVal

	Running int32
	HasPut  int32
	Wg      sync.WaitGroup
}

func (p *JobProc) Val() (string, error) {
	b, err := json.Marshal(&p.JobProcVal)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
