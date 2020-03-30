package main

import (
	"github.com/keepondream/article/migrations"
	"github.com/keepondream/article/route"
)

func main() {

	// 初始化数据库表结构
	migrations.AutoInit()
	// 启动
	route.GinRun(":8080")
}
