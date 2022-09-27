package models

// 执行 cron cmd 的进程
// 注册到 /cronsun/node/<id>
type Node struct {
	ID       int    `json:"id" gorm:"id"`   // machine id
	PID      string `json:"pid" gorm:"pid"` // 进程 pid
	IP       string `json:"ip" gorm:"ip"`   // node ip
	Hostname string `json:"hostname" gorm:"hostname"`
	UUID     string `json:"uuid" gorm:"uuid"`
	Version  string `json:"version" gorm:"version"`
	UpTime   int64  `json:"up" gorm:"up"`     // 启动时间
	DownTime int64  `json:"down" gorm:"down"` // 上次关闭时间

	Status    int  `son:"status" gorm:"status"` // 是否可用
	Connected bool `json:"connected" gorm:"-"`  // 当 Alived 为 true 时有效，表示心跳是否正常
}

func (n *Node) String() string {
	return "node[" + n.UUID + "] pid[" + n.PID + "]"
}
