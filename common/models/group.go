package models

const (
	CronyGroupTypeUser = 1
	CronyGroupTypeNode = 2
)

// 结点类型分组
// 注册到 /cronsun/group/<id>
type Group struct {
	ID      int    `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	Type    int    `json:"type" gorm:"type"`
	Created int64  `json:"created" gorm:"created"`
	Updated int64  `json:"updated" gorm:"updated"`

	NodeIDs []string `json:"nids" gorm:"-"`
}

type NodeGroup struct {
	ID       int    `json:"id" gorm:"id"`
	NodeUUID string `json:"node_uuid" gorm:"node_uuid"`
	GroupId  int    `json:"group_id" gorm:"group_id"`
}

type UserGroup struct {
	ID      int    `json:"id" gorm:"id"`
	UserId  string `json:"user_id" gorm:"user_id"`
	GroupId int    `json:"group_id" gorm:"group_id"`
}
