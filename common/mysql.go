package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MYSQL_HOST     = "127.0.0.1"
	MYSQL_PORT     = 3306
	MYSQL_USER     = "root"
	MYSQL_PASSWORD = "root"
	MYSQL_DBNAME   = "article"
)

// 建立MySQL数据库连接
func Db() *gorm.DB {
	// 连接数据库
	dbStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DBNAME)
	db, err := gorm.Open("mysql", dbStr)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	return db
}
