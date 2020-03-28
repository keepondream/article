package model

// 定义文章表, 为文章表增加 `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
type Article struct {
	Model
	Title   string `gorm:"type:varchar(30)" form:"title" json:"title" binding:"required"`
	Content string `gorm:"type:text" form:"content" json:"content" binding:"required"`
	IsRead  int
	IsOk    int
}
