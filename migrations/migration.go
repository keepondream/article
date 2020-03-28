package migrations

import (
	"github.com/keepondream/article/common"
	"github.com/keepondream/article/model"
)

// migration 数据库迁移
func AutoInit() {
	db := common.Db()
	// 文章表初始化
	article := &model.Article{}
	db.AutoMigrate(article)

}
