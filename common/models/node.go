package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

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

func (n *Node) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyNodeTableName).Create(n).Error
	if err == nil {
		insertId = n.ID
	}
	return
}

// 更新
func (n *Node) Update() error {
	return dbclient.GetMysqlDB().Table(CronyNodeTableName).Updates(n).Error
}

func (n *Node) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where uuid = ?", CronyNodeTableName), n.UUID).Error
}

func (n *Node) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyNodeTableName).Where("uuid = ? ", n.UUID).First(&n).Error
}
