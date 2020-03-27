package main

import (
	"github.com/keepondream/article/migrations"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/keepondream/article/route"
)

var db *gorm.DB

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

func main() {
	// 连接数据库
	var err error
	db, err = gorm.Open("mysql", "root:W8888888w@/article?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrations.AutoInit(db)

	route.GinRun(":8080")
}
