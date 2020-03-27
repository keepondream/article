package model

import (
	"time"
)

// 基本模型的定义
type Model struct {
	ID        uint       `gorm:"primary_key;AUTO_INCREMENT"` // 指定主键,自增
	CreatedAt time.Time  `gorm:"column:create_at"`           // 设置字段名为蛇形小写
	UpdatedAt time.Time  `gorm:"column:update_at"`           // 设置字段名为蛇形小写
	DeletedAt *time.Time `gorm:"column:delete_at"`           // 设置字段名为蛇形小写
}
