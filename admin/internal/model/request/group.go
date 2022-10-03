package request

type (
	ReqGroupSearch struct {
		PageInfo
		ID   int    `json:"id" form:"id"`
		Name string `json:"name" gorm:"name" `
		//分组类型
		//Type int `json:"type" gorm:"type" binding:"required"`
	}
)
