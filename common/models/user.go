package models

type User struct {
	ID       string `json:"id" gorm:"id"`
	UserName string `json:"username" gorm:"username"`
	Password string `json:"password" gorm:"password"`
	Email    string `json:"email" gorm:"email"`
	Role     int    `json:"role" gorm:"email"`
	Status   int    `json:"status" gorm:"status"`

	Created int64 `json:"created" gorm:"created"`
	Updated int64 `json:"updated" gorm:"updated"`
}
