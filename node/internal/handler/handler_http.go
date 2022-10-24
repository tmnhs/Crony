package handler

import (
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/httpclient"
	"strings"
	"time"
)

type HTTPHandler struct {
}

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
	} else {
		urlFields := strings.Split(job.Command, "?")
		url := urlFields[0]
		var body string
		if len(urlFields) >= 2 {
			body = urlFields[1]
		}
		result, err = httpclient.PostJson(url, body, job.Timeout)
	}
	return
}
