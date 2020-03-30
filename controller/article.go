package controller

import (
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
	// 此处个人理解 避免了获取不同的请求头类型的参数问题,还可以有效的避免错误,则定义了_form进行c.bind操作
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
	var article model.Article
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Failed(c, common.WithMsg("请求参数有误"))
		return
	}
	db := common.GetDB()
	db.Debug().First(&article, id)
	if article.ID == 0 {
		common.Failed(c, common.WithMsg("删除失败"))
		return
	}
	db.Debug().Delete(&article)
	common.Success(c, common.WithMsg("删除成功"))
}

// 列表
func ArticleList(c *gin.Context) {
	var articles []model.Article
	db := common.GetDB()
	db.Find(&articles)
	res := make(map[string]interface{})
	res["list"] = &articles
	common.Success(c, common.WithData(res))
}

// 详情
func ArticleDetail(c *gin.Context) {
	var form model.Article
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		common.Failed(c, common.WithMsg("请求参数有误"))
		return
	}
	db := common.GetDB()
	db.Debug().First(&form, id)

	if form.ID == 0 {
		common.Failed(c, common.WithMsg("数据不存在"))
		return
	}
	common.Success(c, common.WithData(common.StructToMapViaJson(form)))
}
