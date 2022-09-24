package utils

import (
	"github.com/google/uuid"
)

func UUID() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
