package resp

import "github.com/tmnhs/crony/common/models"

type (
	RspNodeSearch struct {
		models.Node
		JobCount int `json:"job_count"`
	}
)
