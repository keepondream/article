package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
	"github.com/keepondream/article/model"
)

// 创建
func ArticleCreate(c *gin.Context) {
	var form model.Article
	if c.Bind(&form) != nil {
		common.Failed(c, common.WithMsg("请求参数有误!"))
	}
	db := common.Db()
	fmt.Println(db.Create(&form))

}

// 更新
func ArticleUpdate(c *gin.Context) {

}

// 删除
func ArticleDelete(c *gin.Context) {

}

// 列表
func ArticleList(c *gin.Context) {

}

// 详情
func ArticleDetail(c *gin.Context) {

}
