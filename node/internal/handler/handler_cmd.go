package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"os/exec"
	"time"
)

type CMDHandler struct {
}

func (c *CMDHandler) Run(job *Job) (result string, err error) {
	var (
		cmd  *exec.Cmd
		proc *JobProc
	)
	if job.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, job.Cmd[0], job.Cmd[1:]...)
	} else {
		cmd = exec.Command(job.Cmd[0], job.Cmd[1:]...)
	}
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	err = cmd.Start()
	result = b.String()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("%s\n%s", b.String(), err.Error()))
		return
	}
	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:       cmd.Process.Pid,
			JobID:    job.ID,
			NodeUUID: job.RunOn,
			JobProcVal: models.JobProcVal{
				Time:   time.Now(),
				Killed: false,
			},
		},
	}
	err = proc.Start()
	if err != nil {
		return
	}
	defer proc.Stop()
	if err = cmd.Wait(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("%s%s", b.String(), err.Error()))
		return b.String(), err
	}
	return b.String(), nil
}

func RunPresetScript(script *models.Script) (result string, err error) {
	var cmd *exec.Cmd
	cmd = exec.Command(script.Cmd[0], script.Cmd[1:]...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	err = cmd.Start()
	result = b.String()
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("run preset script:%s\n%s", b.String(), err.Error()))
		return
	}
	if err = cmd.Wait(); err != nil {
		logger.GetLogger().Error(fmt.Sprintf("run preset script:%s%s", b.String(), err.Error()))
		return b.String(), err
	}
	return b.String(), nil
}
