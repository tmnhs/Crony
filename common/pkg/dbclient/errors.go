package dbclient

import "errors"

var (
	ErrClientNotFound   = errors.New("mysql client not found")
	ErrClientDbNameNull = errors.New("mysql dbname is null")
)
