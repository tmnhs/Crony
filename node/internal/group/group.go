package group

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils/errors"
	"strings"
)

type Group struct {
	*models.Group
}
type Groups map[string]*Group

func WatchGroups() clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdGroup, clientv3.WithPrefix(), clientv3.WithPrevKV())
}

func GetGroupById(groupId string) (group *Group, err error) {
	if len(groupId) == 0 {
		return
	}
	resp, err := etcdclient.Get(etcdclient.KeyEtcdGroup + groupId)
	if err != nil || resp.Count == 0 {
		return
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &group)
	return
}

// GetGroups 获取包含 nodeId 的 group
// 如果 nodeId 为空，则获取所有的 group
func GetGroups(nodeId string) (groups map[string]*Group, err error) {
	resp, err := etcdclient.Get(etcdclient.KeyEtcdGroup, clientv3.WithPrefix())
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
			logger.Warnf("group[%s] umarshal err: %s", string(g.Key), e.Error())
			continue
		}
		//获取全部group或者包含nodeId的group
		if len(nodeId) == 0 || group.Included(nodeId) {
			groups[group.ID] = group
		}
	}
	return
}

func GetGroupFromKv(key, value []byte) (g *Group, err error) {
	g = new(Group)
	if err = json.Unmarshal(value, g); err != nil {
		err = fmt.Errorf("group[%s] umarshal err: %s", string(key), err.Error())
	}
	return
}

func DeleteGroupById(id string) (*clientv3.DeleteResponse, error) {
	return etcdclient.Delete(GetGroupKey(id))
}

func GetGroupKey(groupId string) string {
	return etcdclient.KeyEtcdGroup + groupId
}

//func (g *Group) Key() string {
//	return GroupKey(g.ID)
//}
//
func (group *Group) Put(modRev int64) (*clientv3.PutResponse, error) {
	b, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	return etcdclient.PutWithModRev(GetGroupKey(group.ID), string(b), modRev)
}
func IsValidAsKeyPath(s string) bool {
	return strings.IndexAny(s, "/\\") == -1
}
func (group *Group) Check() error {
	group.ID = strings.TrimSpace(group.ID)
	if !IsValidAsKeyPath(group.ID) {
		return errors.ErrIllegalNodeGroupId
	}

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
