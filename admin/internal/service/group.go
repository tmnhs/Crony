package service

import (
	"fmt"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type Groups map[int]*models.Group

//todo
func GetGroupById(groupId int) (group *models.Group, err error) {
	if groupId <= 0 {
		return
	}
	group = &models.Group{}
	err = dbclient.GetMysqlDB().Table(models.CronyGroupTableName).Where("id = ?", groupId).First(group).Error
	return
}

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
