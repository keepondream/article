package main

import (
	"github.com/keepondream/article/migrations"
	"github.com/keepondream/article/route"
)

func main() {

	// client := common.GetRedis()
	// // 第三个参数代表key的过期时间，0代表不会过期。
	// err := client.Set("key", "value2222", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	// val, err := client.Get("key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// db := common.GetDB()

	// fmt.Println(db, client)

	// fmt.Println(val)

	// 初始化数据库表结构
	migrations.AutoInit()
	// 启动
	route.GinRun(":8080")
}
