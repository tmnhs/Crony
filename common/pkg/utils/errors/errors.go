package errors

import "errors"

var (
	ErrClientNotFound   = errors.New("mysql client not found")
	ErrClientDbNameNull = errors.New("mysql dbname is null")
	ErrValueMayChanged  = errors.New("The value has been changed by others on this time.")
	ErrEtcdNotInit      = errors.New("etcd is not initialized")

	ErrNotFound = errors.New("Record not found.")

	ErrEmptyJobName        = errors.New("Name of job is empty.")
	ErrEmptyJobCommand     = errors.New("Command of job is empty.")
	ErrIllegalJobId        = errors.New("Invalid id that includes illegal characters such as '/' '\\'.")
	ErrIllegalJobGroupName = errors.New("Invalid job group name that includes illegal characters such as '/' '\\'.")

	ErrEmptyScriptName    = errors.New("Name of script is empty.")
	ErrEmptyScriptCommand = errors.New("Command of script is empty.")
	ErrEmptyNodeGroupName = errors.New("Name of node group is empty.")
	ErrIllegalNodeGroupId = errors.New("Invalid node group id that includes illegal characters such as '/'.")
)
