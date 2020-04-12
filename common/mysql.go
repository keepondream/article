package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 定义MySQL连接参数常量
const (
	MYSQL_HOST = "127.0.0.1"
	MYSQL_PORT = 33065
	MYSQL_USER = "root"
	// MYSQL_PASSWORD = "W8888888w"
	MYSQL_PASSWORD = "root"
	MYSQL_DBNAME   = "article"
	REDIS_ADDR     = "127.0.0.1"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
	REDIS_DB       = 0
)

var _db *gorm.DB
var err error

// 初始化函数,golang特性,每个包初始化的时候会自动执行init函数,这里用来初始化gorm
func init() {
	// 拼接dsn参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DBNAME)
	fmt.Println(dsn)
	// 连接MySQL,获取DB类型实例,用于后面的数据库读写操作
	_db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic("连接数据库失败,error=" + err.Error())
	}
	// 这里不启用延时关闭数据库连接,因为需要使用连接池,否则会报错
	// defer _db.Close()

	// 设置数据库连接池参数
	_db.DB().SetMaxOpenConns(100) // 设置数据库连接池最大连接数
	_db.DB().SetMaxIdleConns(20)  // 连接池最大允许的空闲连接数,如果没有SQL任务需要执行的连接数大于20,超过的连接会被连接池关闭

	_db.LogMode(true)

}

// 获取gorm db对象, 其他包需要执行数据库查询的时候,只需要通过common.GetDB()获取db对象接口
// 不用担心协程并发使用同样的db对象会共用同一个连接,db对象在调用他的方法的时候会冲数据库连接池中获取新的连接
func GetDB() *gorm.DB {
	return _db
}
