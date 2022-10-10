package handler

import (
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/httpclient"
	"strings"
	"time"
)

type HTTPHandler struct {
}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(job *Job) (result string, err error) {
	var (
		proc *JobProc
	)
	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:       0,
			JobID:    job.ID,
			NodeUUID: job.RunOn,
			JobProcVal: models.JobProcVal{
				Time: time.Now(),
			},
		},
	}

	err = proc.Start()
	if err != nil {
		return
	}
	defer proc.Stop()
	if job.Timeout <= 0 || job.Timeout > HttpExecTimeout {
		job.Timeout = HttpExecTimeout
	}
	if job.HttpMethod == models.HTTPMethodGet {
		result, err = httpclient.Get(job.Command, job.Timeout)
	} else if job.HttpMethod == models.HTTPMethodPost {
		urlFields := strings.Split(job.Command, "?")
		job.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		result, err = httpclient.PostParams(job.Command, params, job.Timeout)
	}
	return
}
