package request

type (
	// PageInfo Paging common input parameter structure
	PageInfo struct {
		Page     int `json:"page" form:"page"`           // 页码
		PageSize int `json:"page_size" form:"page_size"` // 每页大小
	}
	ByID struct {
		ID int `json:"id" form:"id"`
	}
	ByIDS struct {
		IDS []int `json:"ids" form:"ids"`
	}
)

func (page *PageInfo) Check() {
	if page.PageSize <= 0 {
		page.PageSize = 20
	}
	if page.Page <= 0 {
		page.Page = 1
	}
}
