package dbclient

import "errors"

var (
	ErrorProviderNotInit = errors.New("provider not initialized")
	ErrClientNotFound = errors.New("mysql client not found")
)