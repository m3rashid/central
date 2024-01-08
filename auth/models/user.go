package models

const USER_TABLE_NAME = "users"

type User struct {
	BaseModel
	Name         string `json:"name" gorm:"column:name" validate:"required"`
	Email        string `json:"email" gorm:"column:email,unique" validate:"required"`
	Password     string `json:"password" gorm:"column:password" validate:"required"`
	OTP          string `json:"otp" gorm:"column:otp" validate:"required"`
	UserVerified bool   `json:"user_verified" gorm:"column:user_verified,default:false" validate:""`
}

func (*User) TableName() string {
	return USER_TABLE_NAME
}
