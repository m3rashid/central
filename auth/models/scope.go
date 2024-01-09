package models

const SCOPE_TABLE_NAME = "scopes"

type Scope struct {
	BaseModel
	Name        string        `json:"name" gorm:"column:name" validate:"required"`
	Description string        `json:"description" gorm:"column:description" validate:""`
	Permissions []*Permission `json:"permissions" gorm:"many2many:permissions_scope_relation" validate:""`
}

func (*Scope) TableName() string {
	return SCOPE_TABLE_NAME
}

const PERMISSION_TABLE_NAME = "permissions"

type Permission struct {
	BaseModel
	Name string `json:"name" gorm:"column:name" validate:"required"`
}

func (*Permission) TableName() string {
	return PERMISSION_TABLE_NAME
}
