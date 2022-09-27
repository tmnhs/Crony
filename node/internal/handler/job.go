package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jakecoffman/cron"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/common/pkg/utils/errors"
	"strings"
	"time"
)

type Job struct {
	*models.Job
}
type Jobs map[string]*Job

func JobKey(nodeId string, group, id string) string {
	return etcdclient.KeyEtcdJob + nodeId + "/" + group + "/" + id
}

// Note: this function did't check the job.
func GetJob(nodeId, group, id string) (job *Job, err error) {
	job, _, err = GetJobAndRev(nodeId, group, id)
	return
}

func (j *Job) alone() {
	if j.Kind == models.KindAlone {
		j.Parallels = 1
	}
}

func (j *Job) splitCmd() {
	ps := strings.SplitN(j.Command, " ", 2)
	if len(ps) == 1 {
		j.Cmd = ps
		return
	}

	j.Cmd = make([]string, 0, 2)
	j.Cmd = append(j.Cmd, ps[0])
	j.Cmd = append(j.Cmd, utils.ParseCmdArguments(ps[1])...)
}

func (j *Job) String() string {
	data, err := json.Marshal(j)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func GetJobAndRev(nodeId, group, id string) (job *Job, rev int64, err error) {
	resp, err := etcdclient.Get(JobKey(nodeId, group, id))
	if err != nil {
		return
	}

	if resp.Count == 0 {
		err = errors.ErrNotFound
		return
	}

	rev = resp.Kvs[0].ModRevision
	if err = json.Unmarshal(resp.Kvs[0].Value, &job); err != nil {
		return
	}

	job.splitCmd()
	return
}

func DeleteJob(nodeId, group, id string) (resp *clientv3.DeleteResponse, err error) {
	return etcdclient.Delete(JobKey(nodeId, group, id))
}

func GetJobs(nodeId string) (jobs map[string]*Job, err error) {
	resp, err := etcdclient.Get(etcdclient.KeyEtcdJob+nodeId, clientv3.WithPrefix())
	if err != nil {
		return
	}

	count := len(resp.Kvs)
	jobs = make(map[string]*Job, count)
	if count == 0 {
		return
	}

	for _, j := range resp.Kvs {
		job := new(Job)
		if e := json.Unmarshal(j.Value, job); e != nil {
			logger.Warnf("job[%s] umarshal err: %s", string(j.Key), e.Error())
			continue
		}
		//todo
		if err := job.Valid(); err != nil {
			logger.Warnf("job[%s] is invalid: %s", string(j.Key), err.Error())
			continue
		}
		//todo 执行类型
		//job.alone()
		jobs[job.ID] = job
	}
	return
}

func CreateJob(j *Job) cron.FuncJob {
	h := CreateHandler(j)
	if h == nil {
		//logger and error
		return nil
	}
	taskFunc := func() {
		/*handler.taskCount.Add()
		defer handler.taskCount.Done()

		handler.concurrencyQueue.Add()
		defer handler.concurrencyQueue.Done()*/
		logger.Infof("开始执行任务#%s#命令-%s", j.Name, j.Command)
		// 默认只运行任务一次
		var execTimes int = 1
		if j.RetryTimes > 0 {
			execTimes += j.RetryTimes
		}
		var i int = 0
		var output string
		var err error
		for i < execTimes {
			output, err = h.Run(j)
			if err == nil {
				//执行成功
				//todo insert into db
				return
			}
			i++
			if i < execTimes {
				logger.Warnf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", j.ID, i, output, err.Error())
				if j.RetryInterval > 0 {
					time.Sleep(time.Duration(j.RetryInterval) * time.Second)
				} else {
					// 默认重试间隔时间，每次递增1分钟
					time.Sleep(time.Duration(i) * time.Minute)
				}
			}
		}
		// todo 任务执行后置操作 发邮件等
		logger.Infof("任务完成#%s#命令-%s", j.Name, j.Command)
	}
	return taskFunc
}

func (j *Job) Check() error {
	j.ID = strings.TrimSpace(j.ID)
	//todo
	if !IsValidAsKeyPath(j.ID) {
		return errors.ErrIllegalJobId
	}

	j.Name = strings.TrimSpace(j.Name)
	if len(j.Name) == 0 {
		return errors.ErrEmptyJobName
	}
	//todo
	//j.Group = strings.TrimSpace(j.Group)
	//if len(j.Group) == 0 {
	//	j.Group = DefaultJobGroup
	//}

	if !IsValidAsKeyPath(j.Group) {
		return errors.ErrIllegalJobGroupName
	}

	if j.LogExpiration < 0 {
		j.LogExpiration = 0
	}

	j.User = strings.TrimSpace(j.User)

	// 不修改 Command 的内容，简单判断是否为空
	if len(strings.TrimSpace(j.Command)) == 0 {
		return errors.ErrEmptyJobCommand
	}

	return j.Valid()
}
func WatchJobs(nodeId string) clientv3.WatchChan {
	return etcdclient.Watch(etcdclient.KeyEtcdJob+nodeId, clientv3.WithPrefix())
}

func GetJobFromKv(key, value []byte) (job *Job, err error) {
	job = new(Job)
	if err = json.Unmarshal(value, job); err != nil {
		err = fmt.Errorf("job[%s] umarshal err: %s", string(key), err.Error())
		return
	}
	//todo
	//err = job.Valid()
	//job.alone()
	return
}

/*func IsValidAsKeyPath(s string) bool {
	return strings.IndexAny(s, "/\\") == -1
}*/

// 安全选项验证
func (j *Job) Valid() error {
	if len(j.Cmd) == 0 {
		j.splitCmd()
	}
	//todo
	//if err := j.ValidRules(); err != nil {
	//	return err
	//}

	//security := conf.Config.Security
	//if !security.Open {
	//	return nil
	//}
	//
	//if !j.validUser() {
	//	return ErrSecurityInvalidUser
	//}
	//
	//if !j.validCmd() {
	//	return ErrSecurityInvalidCmd
	//}

	return nil
}

func ModifyJob(job *Job) {
	//todo into db

	/*	n.link.delJob(oJob)
		prevCmds := oJob.Cmds(n.ID, n.groups)

		job.Count = oJob.Count
		*oJob = *job
		cmds := oJob.Cmds(n.ID, n.groups)

		for id, cmd := range cmds {
			n.modCmd(cmd, true)
			delete(prevCmds, id)
		}

		for _, cmd := range prevCmds {
			n.delCmd(cmd)
		}

		n.link.addJob(oJob)*/
}

// 从 job etcd 的 key 中取 id
func GetJobIDFromKey(key string) string {
	index := strings.LastIndex(key, "/")
	if index < 0 {
		return ""
	}
	return key[index+1:]
}

//todo db

func (j *Job) Insert2Db() {
	//dbclient.GetMysqlDB()
}
