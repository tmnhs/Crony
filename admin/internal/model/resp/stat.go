package resp

type (
	RspSystemStatistics struct {
		NormalNodeCount    int64 `json:"normal_node_count"`     //正常节点数量
		FailNodeCount      int64 `json:"fail_node_count"`       //不正常节点数量
		JobExcSuccessCount int64 `json:"job_exc_success_count"` //任务执行总数
		JobRunningCount    int64 `json:"job_running_count"`     //任务正在执行总数
		JobExcFailCount    int64 `json:"job_exc_fail_count"`    //任务执行失败总数
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
