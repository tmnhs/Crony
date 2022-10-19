package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

const (
	RoleNormal = 1
	RoleAdmin  = 2
)

type User struct {
	ID       int    `json:"id" gorm:"column:id;primary_key;auto_increment"`
	UserName string `json:"username" gorm:"size:128;column:username;not null"`
	Password string `json:"password" gorm:"size:128;column:password;not null"`
	Email    string `json:"email" gorm:"size:64;column:email;default:''"`
	Role     int    `json:"role" gorm:"size:1;column:role;default:1"`

	Created int64 `json:"created" gorm:"column:created;not null"`
	Updated int64 `json:"updated" gorm:"column:updated;default:0"`
}

func (u *User) Update() error {
	return dbclient.GetMysqlDB().Table(CronyUserTableName).Updates(u).Error
}

func (u *User) Delete() error {
	return dbclient.GetMysqlDB().Exec(fmt.Sprintf("delete from %s where id = ?", CronyUserTableName), u.ID).Error
}

func (u *User) Insert() (insertId int, err error) {
	err = dbclient.GetMysqlDB().Table(CronyUserTableName).Create(u).Error
	if err == nil {
		insertId = u.ID
	}
	return
}

func (u *User) FindById() error {
	return dbclient.GetMysqlDB().Table(CronyUserTableName).Select("id", "username", "email", "role", "created", "updated").Where("id = ? ", u.ID).First(u).Error
}

func (u *User) TableName() string {
	return CronyUserTableName
}
