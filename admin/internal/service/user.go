package service

import (
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
	"github.com/tmnhs/crony/common/pkg/utils"
)

type UserService struct {
}

var DefaultUserService = new(UserService)

func (us *UserService) Login(username, password string) (u *models.User, err error) {
	err = dbclient.GetMysqlDB().Table(models.CronyUserTableName).Where("username = ? And password = ?", username, utils.MD5(password)).Find(u).Error
	return
}

func (us *UserService) FindByUserName(username string) (u *models.User, err error) {
	err = dbclient.GetMysqlDB().Table(models.CronyUserTableName).Where("username = ? ", username).Find(u).Error
	return
}

func (us *UserService) ChangePassword(userId int, oldPassword, newPassword string) error {
	return dbclient.GetMysqlDB().Table(models.CronyUserTableName).Where("id = ? And password ", userId, utils.MD5(oldPassword)).Update("password", utils.MD5(newPassword)).Error
}

func (us *UserService) Search(s *request.ReqUserSearch) ([]models.User, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyUserTableName)
	if len(s.UserName) > 0 {
		db = db.Where("username like ?", s.UserName+"%")
	}
	if len(s.Email) > 0 {
		db.Where("email = ?", s.Email)
	}
	if s.Role > 0 {
		db.Where("role = ?", s.Role)
	}
	users := make([]models.User, 2)
	var total int64
	err := db.Select("id username email role").Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
