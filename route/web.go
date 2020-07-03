package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/controller"
)

// Web web端接口
func Web(r *gin.Engine, prefix string) {
	web := r.Group(prefix)

	web.GET("test", func(c *gin.Context) {
		fmt.Println("web test success")
	})

	web.Any("eolinker", controller.Eolinker)
}
