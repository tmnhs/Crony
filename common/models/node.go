package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

const (
	NodeConnSuccess = 1
	NodeConnFail    = 2

	NodeSystemInfoSwitch = "alive"
)

// 注册到 /crony/node/<node_uuid>>/
type Node struct {
	ID       int    `json:"id" gorm:"column:id"`   // machine id
	PID      string `json:"pid" gorm:"column:pid"` // 进程 pid
	IP       string `json:"ip" gorm:"column:ip"`   // node ip
	Hostname string `json:"hostname" gorm:"column:hostname"`
	UUID     string `json:"uuid" gorm:"column:uuid"`
	Version  string `json:"version" gorm:"column:version"`
	UpTime   int64  `json:"up" gorm:"column:up"`     // 启动时间
	DownTime int64  `json:"down" gorm:"column:down"` // 上次关闭时间

	Status int `json:"status" gorm:"column:status"` // 是否可用
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

func (n *Node) FindByUUID() error {
	return dbclient.GetMysqlDB().Table(CronyNodeTableName).Where("uuid = ? ", n.UUID).First(n).Error
}
