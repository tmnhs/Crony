package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

const (
	CronyGroupTypeUser = 1
	CronyGroupTypeNode = 2
)

type Group struct {
	ID   int    `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name" binding:"required"`
	//分组类型
	Type    int   `json:"type" gorm:"type" binding:"required"`
	Created int64 `json:"created" gorm:"created"`
	Updated int64 `json:"updated" gorm:"updated"`

	NodeIDs []string `json:"nids" gorm:"-"`
}

func (g *Group) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *Group) Update() error {
	return dbclient.GetMysqlDB().Table(CronyGroupTableName).Updates(g).Error
}
func (g *Group) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where id = ?", CronyGroupTableName), g.ID).Error
}
func (g *Group) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyGroupTableName).Where("id = ? ", g.ID).First(g).Error
}

type NodeGroup struct {
	ID       int    `json:"id" gorm:"id"`
	NodeUUID string `json:"node_uuid" gorm:"node_uuid" binding:"required"`
	GroupId  int    `json:"group_id" gorm:"group_id" binding:"required"`
}

func (g *NodeGroup) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyNodeGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *NodeGroup) Update() error {
	return dbclient.GetMysqlDB().Table(CronyNodeGroupTableName).Updates(g).Error
}
func (g *NodeGroup) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where node_uuid = ? and group_id = ?", CronyNodeGroupTableName), g.NodeUUID, g.GroupId).Error
}
func (g *NodeGroup) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyNodeGroupTableName).Where("id = ? ", g.ID).First(g).Error
}

type UserGroup struct {
	ID      int `json:"id" gorm:"id"`
	UserId  int `json:"user_id" gorm:"user_id" binding:"required"`
	GroupId int `json:"group_id" gorm:"group_id" binding:"required" `
}

func (g *UserGroup) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyUserGroupTableName).Create(g).Error
	if err == nil {
		insertId = g.ID
	}
	return
}
func (g *UserGroup) Update() error {
	return dbclient.GetMysqlDB().Table(CronyUserGroupTableName).Updates(g).Error
}
func (g *UserGroup) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where group_id = ? and  user_id =?", CronyUserGroupTableName), g.GroupId, g.UserId).Error
}
func (g *UserGroup) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyUserGroupTableName).Where("id = ? ", g.ID).First(g).Error
}
