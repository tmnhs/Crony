package handler

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jakecoffman/cron"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/config"
	"github.com/tmnhs/crony/common/pkg/etcdclient"
	"github.com/tmnhs/crony/common/pkg/logger"
	"github.com/tmnhs/crony/common/pkg/notify"
	"github.com/tmnhs/crony/common/pkg/utils"
	"github.com/tmnhs/crony/common/pkg/utils/errors"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	*models.Job
}
type Jobs map[int]*Job

func JobKey(nodeUUID string, jobId int) string {
	return fmt.Sprintf(etcdclient.KeyEtcdJob, nodeUUID, jobId)
}

// Note: this function did't check the job.
func GetJob(nodeUUID string, jobId int) (job *Job, err error) {
	job, _, err = GetJobAndRev(nodeUUID, jobId)
	return
}

/*func (j *Job) alone() {
	if j.Kind == models.KindAlone {
		j.Parallels = 1
	}
}*/

func (j *Job) String() string {
	data, err := json.Marshal(j)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func GetJobAndRev(nodeUUID string, jobId int) (job *Job, rev int64, err error) {
	resp, err := etcdclient.Get(JobKey(nodeUUID, jobId))
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

	job.SplitCmd()
	return
}

func DeleteJob(nodeUUID string, jobId int) (resp *clientv3.DeleteResponse, err error) {
	return etcdclient.Delete(JobKey(nodeUUID, jobId))
}

func GetJobs(nodeUUID string) (jobs Jobs, err error) {
	resp, err := etcdclient.Get(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
	if err != nil {
		return
	}

	count := len(resp.Kvs)
	jobs = make(Jobs, count)
	if count == 0 {
		return
	}

	for _, j := range resp.Kvs {
		job := new(Job)
		if e := json.Unmarshal(j.Value, job); e != nil {
			logger.GetLogger().Warn(fmt.Sprintf("job[%s] umarshal err: %s", string(j.Key), e.Error()))
			continue
		}
		if err := job.Check(); err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("job[%s] is invalid: %s", string(j.Key), err.Error()))
			continue
		}
		jobs[job.ID] = job
	}
	return
}

func (j *Job) RunWithRecovery() {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logger.GetLogger().Warn(fmt.Sprintf("panic running job: %v\n%s", r, buf))
		}
	}()
	t := time.Now()
	jobLogId, err := j.CreateJobLog()
	if err != nil {
		logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
	}
	h := CreateHandler(j)
	if h == nil {
		//logger and error
		return
	}
	result, err := h.Run(j)
	if err != nil {
		err = j.Fail(jobLogId, t, err.Error(), 0)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
		}
	} else {
		err = j.Success(jobLogId, t, result, 0)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
		}
	}
}

func CreateJob(j *Job) cron.FuncJob {
	h := CreateHandler(j)
	if h == nil {
		//logger and error
		return nil
	}
	jobFunc := func() {
		/*handler.taskCount.Add()
		defer handler.taskCount.Done()

		handler.concurrencyQueue.Add()
		defer handler.concurrencyQueue.Done()*/
		logger.GetLogger().Info(fmt.Sprintf("start the job#%s#command-%s", j.Name, j.Command))
		// 默认只运行任务一次
		var execTimes int = 1
		if j.RetryTimes > 0 {
			execTimes += j.RetryTimes
		}
		var i = 0
		var output string
		var runErr error
		var err error
		var jobLogId int
		t := time.Now()
		jobLogId, err = j.CreateJobLog()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
		}
		for i < execTimes {
			output, runErr = h.Run(j)
			if runErr == nil {
				//执行成功
				err = j.Success(jobLogId, t, output, i)
				if err != nil {
					logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
				}
				return
			}
			i++
			if i < execTimes {
				logger.GetLogger().Warn(fmt.Sprintf("job execution failure#jobId-%d#retry%d次#output-%s#error-%s", j.ID, i, output, err.Error()))
				if j.RetryInterval > 0 {
					time.Sleep(time.Duration(j.RetryInterval) * time.Second)
				} else {
					// 默认重试间隔时间，每次递增1分钟
					time.Sleep(time.Duration(i) * time.Minute)
				}
			}
		}
		//执行全部失败
		err = j.Fail(jobLogId, t, err.Error(), execTimes-1)
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to write to job log with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
		}
		node := &models.Node{UUID: j.RunOn}
		err = node.FindByUUID()
		if err != nil {
			logger.GetLogger().Warn(fmt.Sprintf("Failed to find node with jobID:%d nodeUUID: %s error:%s", j.ID, j.RunOn, err.Error()))
		}
		var to []string
		for _, userId := range j.NotifyToArray {
			userModel := &models.User{ID: userId}
			err = userModel.FindById()
			if err != nil {
				continue
			}
			if j.NotifyType == notify.NotifyTypeMail {
				to = append(to, userModel.Email)
			} else if j.NotifyType == notify.NotifyTypeMail && config.GetConfigModels().WebHook.Kind == "feishu" {
				to = append(to, userModel.UserName)
			}
		}
		msg := &notify.Message{
			Type:      j.NotifyType,
			IP:        fmt.Sprintf("%s:%s", node.IP, node.PID),
			Subject:   "",
			Body:      fmt.Sprintf("job[%d] run on node[%s] execute failed ,error info :%s", j.ID, j.RunOn, runErr.Error()),
			To:        to,
			OccurTime: time.Now().Format(utils.TimeFormatSecond),
		}
		go notify.Send(msg)
	}
	return jobFunc
}
func WatchJobs(nodeUUID string) clientv3.WatchChan {
	return etcdclient.Watch(fmt.Sprintf(etcdclient.KeyEtcdJobProfile, nodeUUID), clientv3.WithPrefix())
}

func GetJobFromKv(key, value []byte) (job *Job, err error) {
	job = new(Job)
	if err = json.Unmarshal(value, job); err != nil {
		err = fmt.Errorf("job[%s] umarshal err: %s", string(key), err.Error())
		return
	}
	err = job.Check()
	return
}

// 从 etcd 的 key 中取 job_id
func GetJobIDFromKey(key string) int {
	index := strings.LastIndex(key, "/")
	if index < 0 {
		return 0
	}
	jobId, err := strconv.Atoi(key[index+1:])
	if err != nil {
		return 0
	}
	return jobId
}

//将每次执行任务的结果写入日志
func (j *Job) CreateJobLog() (int, error) {
	start := time.Now()
	jobLog := &models.JobLog{
		Name:      j.Name,
		JobId:     j.ID,
		Command:   j.Command,
		IP:        j.Ip,
		Hostname:  j.Hostname,
		NodeUUID:  j.RunOn,
		Spec:      j.Spec,
		StartTime: start.Unix(),
	}
	return jobLog.Insert()
}

func UpdateJobLog(jobLogId int, start time.Time, output string, retry int, success bool) error {
	end := time.Now()
	jobLog := &models.JobLog{
		ID:         jobLogId,
		StartTime:  start.Unix(),
		RetryTimes: retry,
		Success:    success,
		Output:     output,
		EndTime:    end.Unix(),
	}
	return jobLog.Update()
}
func (j *Job) Success(jobLogId int, start time.Time, output string, retry int) error {
	return UpdateJobLog(jobLogId, start, output, retry, true)
}

func (j *Job) Fail(jobLogId int, start time.Time, errMsg string, retry int) error {
	return UpdateJobLog(jobLogId, start, errMsg, retry, false)
}
