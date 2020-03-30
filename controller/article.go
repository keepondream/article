package controller

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
	"github.com/keepondream/article/model"
)

// 创建
func ArticleCreate(c *gin.Context) {
	var form model.Article
	if c.Bind(&form) != nil {
		common.Failed(c, common.WithMsg("请求参数有误!"))
		return
	}
	db := common.GetDB()
	db.Save(&form)
	if form.ID > 0 {
		common.Success(c, common.WithData(common.StructToMapViaJson(form)))
	} else {
		common.Failed(c, common.WithMsg("文章创建失败"))
		return
	}

}

// 更新
func ArticleUpdate(c *gin.Context) {
	var form model.Article
	var _form model.Article
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Failed(c, common.WithMsg("请求参数有误"))
		return
	}
	c.Bind(&_form)
	db := common.GetDB()
	db.Debug().First(&form, id)

	form.Title = _form.Title
	form.Content = _form.Content
	db.Debug().Save(&form)
	common.Success(c)
}

// 删除
func ArticleDelete(c *gin.Context) {

}

// 列表
func ArticleList(c *gin.Context) {
}

// 详情
func ArticleDetail(c *gin.Context) {
	var form model.Article
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Failed(c, common.WithMsg("请求参数有误"))
		return
	}
	fmt.Println(id)
	db := common.GetDB()
	db.Debug().First(&form, id)

	if form.ID == 0 {
		common.Failed(c, common.WithMsg("数据不存在"))
		return
	}
	common.Success(c, common.WithData(common.StructToMapViaJson(form)))
}
