package etcdclient

import "errors"

var (
	ErrValueMayChanged = errors.New("The value has been changed by others on this time.")
)
