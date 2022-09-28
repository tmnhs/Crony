package handler

import (
	"fmt"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/httpclient"
	"net/http"
	"strings"
)

type HTTPHandler struct {
}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(job *Job) (result string, err error) {
	if job.Timeout <= 0 || job.Timeout > HttpExecTimeout {
		job.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if job.HttpMethod == models.HTTPMethodGet {
		resp = httpclient.Get(job.Command, job.Timeout)
	} else if job.HttpMethod == models.HTTPMethodPost {
		urlFields := strings.Split(job.Command, "?")
		job.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(job.Command, params, job.Timeout)
	}
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}

	return resp.Body, err
}
