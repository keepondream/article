package model

import (
	"github.com/keepondream/article/common"
)

// 基本模型的定义
type Model struct {
	ID        uint             `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`     // 指定主键,自增
	CreatedAt common.JSONTime  `gorm:"column:create_at" form:"create_at" json:"create_at"` // 设置字段名
	UpdatedAt common.JSONTime  `gorm:"column:update_at" form:"update_at" json:"update_at"` // 设置字段名
	DeletedAt *common.JSONTime `gorm:"column:delete_at" form:"delete_at" json:"delete_at"` // 设置字段名
}
