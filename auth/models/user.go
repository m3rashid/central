package models

import (
	internal "internal/models"
)

const USER_TABLE_NAME = "users"

type User struct {
	internal.BaseModel
	Name          string    `json:"name" gorm:"column:name;not null" validate:"required"`
	Email         string    `json:"email" gorm:"column:email;unique;not null" validate:"required"`
	Password      string    `json:"password" gorm:"column:password;not null" validate:"required"`
	OTP           string    `json:"otp" gorm:"column:otp;not null" validate:"required"`
	UserVerified  bool      `json:"user_verified" gorm:"column:user_verified;default:false" validate:""`
	ConnectedApps []*Client `json:"connected_apps" gorm:"many2many:connected_apps_user_relation" validate:""`
	// LastLogin time.Time `json:"last_login" gorm:"column:last_login"` // handle this for password revalidation
}

func (*User) TableName() string {
	return USER_TABLE_NAME
}
