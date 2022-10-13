package request

type (
	ReqNodeSearch struct {
		PageInfo
		IP     string `json:"ip" form:"ip"` // node ip
		UUID   string `json:"uuid" form:"uuid"`
		UpTime int64  `json:"up" form:"up"`        // 启动时间
		Status int    `son:"status" form:"status"` // 状态
	}
	ByUUID struct {
		UUID string `json:"uuid" form:"uuid"`
	}
)
