package service

import (
	"fmt"
	"github.com/tmnhs/crony/admin/internal/model/request"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type Groups map[int]*models.Group

type GroupService struct {
}

var DefaultGroupService = new(GroupService)

// GetGroups 获取包含 nodeId 的 group
// 如果 nodeId 为空，则获取所有的 group
func GetGroups(nodeUUID string) (groupsMap Groups, err error) {
	sql := fmt.Sprintf("select g.id as id g.name as name from %s  ng left join %s g  on  ng.group_id = g.id ")
	if len(nodeUUID) > 0 {
		sql += "and ng.node_uuid=?"
	}
	groups := make([]*models.Group, 2)
	groupsMap = make(Groups, 2)
	err = dbclient.GetMysqlDB().Raw(fmt.Sprintf(sql, models.CronyNodeGroupTableName, models.CronyGroupTableName), nodeUUID).Scan(groups).Error
	if err != nil {
		return
	}
	for _, group := range groups {
		groupsMap[group.ID] = group
	}
	return
}

func (n *GroupService) Search(s *request.ReqGroupSearch) ([]models.Group, int64, error) {
	db := dbclient.GetMysqlDB().Table(models.CronyGroupTableName)
	if len(s.Name) > 0 {
		db = db.Where("name = ?", s.Name)
	}
	if s.Type > 0 {
		db.Where("type = ?", s.Type)
	}
	if s.ID > 0 {
		db.Where("id = ?", s.ID)
	}
	groups := make([]models.Group, 2)
	var total int64
	err := db.Limit(s.PageSize).Offset((s.Page - 1) * s.PageSize).Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return groups, total, nil
}
