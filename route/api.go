package route

import (
	"fmt"

	"github.com/keepondream/article/controller"

	"github.com/gin-gonic/gin"
)

// 获取gin路由引擎
func GinRun(port string) *gin.Engine {
	// 初始化引擎
	route := gin.Default()

	// 设置api路由组
	Api(route, "api")

	// 设置web路由组
	Web(route, "web")

	// 启动
	route.Run(port)

	return route
}

// api接口
func Api(r *gin.Engine, prefix string) {
	api := r.Group(prefix)

	api.GET("login", func(c *gin.Context) {
		fmt.Println("api login success")
	})

	api.GET("test", controller.Test)

	// 文章路由
	api.POST("article", controller.ArticleCreate)
	api.PUT("article", controller.ArticleUpdate)
	api.DELETE("article", controller.ArticleDelete)
	api.GET("article", controller.ArticleDetail)
	api.GET("articleList", controller.ArticleList)
}
