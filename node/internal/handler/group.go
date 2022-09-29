package handler

import (
	"fmt"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/dbclient"
)

type Groups map[int]*models.Group

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

/*
func GetGroupById(gid int) (g *Group, err error) {
	if gid <= 0 {
		return
	}
	resp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdGroup,gid) )
	if err != nil || resp.Count == 0 {
		return
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &g)
	return
}

// GetGroups 获取包含 nid 的 group
// 如果 nid 为空，则获取所有的 group
func GetGroups(nid string) (groups map[string]*Group, err error) {
	resp, err := etcdclient.Get(conf.Config.Group, client.WithPrefix())
	if err != nil {
		return
	}

	count := len(resp.Kvs)
	groups = make(map[string]*Group, count)
	if count == 0 {
		return
	}

	for _, g := range resp.Kvs {
		group := new(Group)
		if e := json.Unmarshal(g.Value, group); e != nil {
			log.Warnf("group[%s] umarshal err: %s", string(g.Key), e.Error())
			continue
		}
		if len(nid) == 0 || group.Included(nid) {
			groups[group.ID] = group
		}
	}
	return
}

func WatchGroups() client.WatchChan {
	return etcdclient.Watch(conf.Config.Group, client.WithPrefix(), client.WithPrevKV())
}

func GetGroupFromKv(key, value []byte) (g *Group, err error) {
	g = new(Group)
	if err = json.Unmarshal(value, g); err != nil {
		err = fmt.Errorf("group[%s] umarshal err: %s", string(key), err.Error())
	}
	return
}

func DeleteGroupById(id string) (*client.DeleteResponse, error) {
	return etcdclient.Delete(GroupKey(id))
}

func GroupKey(id string) string {
	return conf.Config.Group + id
}

func (g *Group) Key() string {
	return GroupKey(g.ID)
}

func (g *Group) Put(modRev int64) (*client.PutResponse, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	return etcdclient.PutWithModRev(g.Key(), string(b), modRev)
}

func (g *Group) Check() error {
	g.ID = strings.TrimSpace(g.ID)
	if !IsValidAsKeyPath(g.ID) {
		return ErrIllegalNodeGroupId
	}

	g.Name = strings.TrimSpace(g.Name)
	if len(g.Name) == 0 {
		return ErrEmptyNodeGroupName
	}

	return nil
}

func (g *Group) Included(nid string) bool {
	for i, count := 0, len(g.NodeIDs); i < count; i++ {
		if nid == g.NodeIDs[i] {
			return true
		}
	}

	return false
}


func (group *Group) Check() error {
	group.Name = strings.TrimSpace(group.Name)
	if len(group.Name) == 0 {
		return errors.ErrEmptyNodeGroupName
	}
	return nil
}
//group是否包含nodeId
func (group *Group) Included(nodeId string) bool {
	for i, count := 0, len(group.NodeIDs); i < count; i++ {
		if nodeId == group.NodeIDs[i] {
			return true
		}
	}
	return false
}
*/
