package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// web端接口
func Web(r *gin.Engine, prefix string) {
	web := r.Group(prefix)

	web.GET("test", func(c *gin.Context) {
		fmt.Println("web test success")
	})
}
