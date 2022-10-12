package handler

import (
	"github.com/tmnhs/crony/common/models"
)

type Handler interface {
	Run(job *Job) (string, error)
}

func CreateHandler(j *Job) Handler {
	var handler Handler = nil
	switch j.Type {
	case models.JobTypeCmd:
		handler = new(CMDHandler)
	case models.JobTypeHttp:
		handler = new(HTTPHandler)
	}
	return handler
}
