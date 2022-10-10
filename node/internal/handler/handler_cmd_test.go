package handler_test

import (
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/node/internal/handler"
	"testing"
)

func TestCmd(t *testing.T) {
	jobs := []handler.Job{
		{
			Job: &models.Job{
				Name:    "test",
				Command: "/home/tmnh/hello.sh",
				Cmd:     []string{"/home/tmnh/hello.sh"},
				Type:    models.JobTypeCmd,
				Timeout: 0,
			},
		},
		//{
		//	Name:          "post",
		//	HttpUrl:       "",
		//	HttpMethod:    models.HTTPMethodPost,
		//	Timeout:       3000,
		//},
	}
	var cmd handler.CMDHandler
	for i := 0; i < len(jobs); i++ {
		rsp, err := cmd.Run(&jobs[i])
		if err != nil {
			t.Error(err)
		}
		t.Log(rsp)
	}
}
