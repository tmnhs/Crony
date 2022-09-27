package models

// 结点类型分组
// 注册到 /cronsun/group/<id>
type Group struct {
	ID      string `json:"id" gorm:"id"`
	Name    string `json:"name" gorm:"name"`
	Type    int    `json:"type" gorm:"type"`
	Created int64  `json:"created" gorm:"created"`
	Updated int64  `json:"updated" gorm:"updated"`

	NodeIDs []string `json:"nids" gorm:"-"`
}
