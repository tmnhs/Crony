package request

type (
	ReqUserLogin struct {
		UserName string `json:"username" form:"username" binding:"required,min=2,max=20"`
		Password string `json:"password" form:"password" binding:"required,min=4,max=20,alphanum"`
	}
	ReqUserRegister struct {
		UserName string `json:"username" form:"username" binding:"required,min=2,max=20"`
		Password string `json:"password" form:"password" binding:"required,min=4,max=20,alphanum"`
		Email    string `json:"email" form:"email"`
		Role     int    `json:"role" form:"email"`
	}
	// Modify password structure
	ReqChangePassword struct {
		Password    string `json:"password" binding:"required,min=4,max=20,alphanum"`
		NewPassword string `json:"new_password" binding:"required,min=4,max=20,alphanum"`
	}
	ReqUserSearch struct {
		PageInfo
		ID       int    `json:"id" form:"id"`
		UserName string `json:"username" form:"username"`
		Email    string `json:"email" form:"email"`
		Role     int    `json:"role" form:"email"`
	}
)
