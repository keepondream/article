package model

// 定义文章表, 为文章表增加 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
type Article struct {
	Model
	Title   string `gorm:"type:varchar(30);default:''" form:"title" json:"title" binding:"required"`
	Content string `gorm:"type:text;" form:"content" json:"content" binding:"required"`
	IsRead  uint8  `gorm:"type:tinyint;default:2;not null"`
	IsOk    uint8  `gorm:"type:tinyint;default:2;not null"`
}
