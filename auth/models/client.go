package models

import (
	internal "internal/models"
)

const CLIENT_TABLE_NAME = "clients"

type Client struct {
	internal.BaseModel
	ClientID           string   `json:"client_id" gorm:"column:client_id" validate:"required"`
	ClientSecret       string   `json:"client_secret" gorm:"column:client_secret" validate:"required"`
	Scopes             []*Scope `json:"scopes" gorm:"many2many:scopes_client_relation" validate:""`
	SuccessRedirectUri string   `json:"success_redirect_uri" gorm:"column:success_redirect_uri" validate:"required"`
	FailureRedirectUri string   `json:"failure_redirect_uri" gorm:"column:failure_redirect_uri" validate:"required"`
	AppName            string   `json:"client_app_name" gorm:"column:client_app_name" validate:""`
	AppLogoUrl         string   `json:"client_app_logo_url" gorm:"column:client_app_logo_url" validate:""`
}

func (*Client) TableName() string {
	return CLIENT_TABLE_NAME
}
