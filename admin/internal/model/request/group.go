package request

type (
	ReqGroupSearch struct {
		PageInfo
		ID   int    `json:"id" form:"id"`
		Name string `json:"name" gorm:"name" `
	}
)
