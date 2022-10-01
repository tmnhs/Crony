package resp

import "github.com/tmnhs/crony/common/models"

type (
	RspLogin struct {
		User  *models.User `json:"user"`
		Token string       `json:"token"`
	}
)
