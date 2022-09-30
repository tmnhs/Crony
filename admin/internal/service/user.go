package service

import (
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type UserService struct {
}

var DefaultUserService = new(UserService)

func (us *UserService) Login(username, password string) (u *models.User, err error) {
	err = dbclient.GetMysqlDB().Table(models.CronyUserTableName).Where("username = ? And password = ?", username, password).Find(u).Error
	return
}
