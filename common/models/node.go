package models

import (
	"time"
)

// 执行 cron cmd 的进程
// 注册到 /cronsun/node/<id>
type Node struct {
	ID       string `bson:"_id" json:"id"`  // machine id
	PID      string `bson:"pid" json:"pid"` // 进程 pid
	IP       string `bson:"ip" json:"ip"`   // node ip
	Hostname string `bson:"hostname" json:"hostname"`

	Version  string    `bson:"version" json:"version"`
	UpTime   time.Time `bson:"up" json:"up"`     // 启动时间
	DownTime time.Time `bson:"down" json:"down"` // 上次关闭时间

	Alived    bool `bson:"alived" json:"alived"` // 是否可用
	Connected bool `bson:"-" json:"connected"`   // 当 Alived 为 true 时有效，表示心跳是否正常
}
