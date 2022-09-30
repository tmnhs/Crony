package models

import (
	"fmt"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type User struct {
	ID       int    `json:"id" gorm:"id"`
	UserName string `json:"username" gorm:"username"`
	Password string `json:"password" gorm:"password"`
	Email    string `json:"email" gorm:"email"`
	Role     int    `json:"role" gorm:"email"`
	Status   int    `json:"status" gorm:"status"`

	Created int64 `json:"created" gorm:"created"`
	Updated int64 `json:"updated" gorm:"updated"`
}

// 更新
func (u *User) Update() error {
	//只会更新非零字段
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
