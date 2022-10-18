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
	ID       int    `json:"id" gorm:"column:id"`
	UserName string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
	Email    string `json:"email" gorm:"column:email"`
	Role     int    `json:"role" gorm:"column:role"`
	Status   int    `json:"status" gorm:"column:status"`

	Created int64 `json:"created" gorm:"column:created"`
	Updated int64 `json:"updated" gorm:"column:updated"`
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
