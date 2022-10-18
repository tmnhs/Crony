package resp

type (
	RspSystemStatistics struct {
		NormalNodeCount    int64 `json:"normal_node_count"`
		FailNodeCount      int64 `json:"fail_node_count"`
		JobExcSuccessCount int64 `json:"job_exc_success_count"`
		JobRunningCount    int64 `json:"job_running_count"`
		JobExcFailCount    int64 `json:"job_exc_fail_count"`
	}
	DateCount struct {
		Date  string `json:"date" gorm:"column:date"`
		Count string `json:"count" gorm:"column:count"`
	}
	RspDateCount struct {
		SuccessDateCount []DateCount `json:"success_date_count"`
		FailDateCount    []DateCount `json:"fail_date_count"`
	}
)
