package migrations

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/keepondream/article/model"
)

func AutoInit(db *gorm.DB) {

	article := &model.Article{}
	db.AutoMigrate(article)
}
