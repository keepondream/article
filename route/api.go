package route

import (
	"fmt"

	"github.com/keepondream/article/common"
	"github.com/keepondream/article/controller"

	"github.com/gin-gonic/gin"
)

// 获取gin路由引擎
func GinRun(port string) *gin.Engine {
	// 初始化引擎
	route := gin.Default()
	// 加载html目录下的所有模板文件
	route.LoadHTMLGlob("html/*")
	// 加载logger 日志记录
	route.Use(common.LoggerToFile())
	// 有子目录，模板文件都在子目录里进行加载
	// route.LoadHTMLGlob("html/**/*")

	// HTML 路劲
	route.GET("/", controller.Html)

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

	api.POST("test", controller.Test)

	// 文章路由
	api.POST("article", controller.ArticleCreate)
	api.PUT("article/:id", controller.ArticleUpdate)
	api.DELETE("article/:id", controller.ArticleDelete)
	api.GET("article/:id", controller.ArticleDetail)
	api.GET("articleList", controller.ArticleList)

	// 上传文件
	api.POST("upload", controller.TransformFile)
	// 下载文件
	api.GET("download/:filename", controller.DownloadFile)

}
