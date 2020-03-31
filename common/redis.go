package common

import "github.com/go-redis/redis"

var _redis *redis.Client

func init() {
	// 初始化一个新的Redis client, go-redis 包自带了连接池,会自动维护Redis连接,因此创建一次client即可,不要查询一次Redis就关闭client
	_redis = redis.NewClient(&redis.Options{
		Addr:     REDIS_ADDR + ":" + REDIS_PORT,
		Password: REDIS_PASSWORD,
		DB:       REDIS_DB,
	})

}

// 获取Redis对象
func GetRedis() *redis.Client {
	return _redis
}
