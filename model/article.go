package model

import "github.com/jinzhu/gorm"

// 定义文章表, 为文章表增加 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
type Article struct {
	gorm.Model
	Title   string `gorm:"type:varchar(30)"`
	Content string
	IsRead  int
}
