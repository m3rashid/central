package models

import (
	internal "internal/models"
	"time"
)

const CONTACT_TABLE_NAME = "contacts"

type Contact struct {
	internal.BaseModel
	UserID   uint      `json:"userId" gorm:"userId;not null" validate:"required"`
	Name     string    `json:"name" gorm:"column:name;not null" validate:"required"`
	Email    string    `json:"email" gorm:"column:email" validate:"email"`
	Phone    string    `json:"phone" gorm:"column:phone" validate:""`
	Birthday time.Time `json:"birthday" gorm:"column:birthday" validate:""`
	Extra    string    `json:"" gorm:"" validate:""` // use JSON type for postgres
}

func (*Contact) TableName() string {
	return CONTACT_TABLE_NAME
}
