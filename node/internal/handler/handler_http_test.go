package handler_test

import (
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/node/internal/handler"
	"testing"
)

func TestHttpCall(t *testing.T) {
	jobs := []handler.Job{
		{
			Job: &models.Job{
				Name:       "get",
				HttpUrl:    "https://www.baidu.com",
				HttpMethod: models.HTTPMethodGet,
				Timeout:    3000,
			},
		},
		// TODO
		//{
		//	Name:          "post",
		//	HttpUrl:       "",
		//	HttpMethod:    models.HTTPMethodPost,
		//	Timeout:       3000,
		//},
	}
	var http handler.HTTPHandler
	for i := 0; i < len(jobs); i++ {
		rsp, err := http.Run(&jobs[i])
		if err != nil {
			t.Error(err)
		}
		t.Log(rsp)
	}
}
