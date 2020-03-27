package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/keepondream/article/common"
)

func Test(c *gin.Context) {
	// common.Success(c)
	common.Success(c, common.WithCode(403), common.WithMsg("坤哥牛逼"))
}
