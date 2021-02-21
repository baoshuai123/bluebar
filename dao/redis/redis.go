package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
)

// 声明一个全局的redisDb变量
var redisDb *redis.Client

// 根据redis配置初始化一个客户端
func Init(cfg *settings.RedisConfig) (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", // redis地址
			//viper.GetString("redis.host"),
			//viper.GetInt("redis.port"),
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password, // redis密码，没有则留空
		DB:       cfg.Db,       // 默认数据库，默认是0
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err = redisDb.Ping().Result()
	return
}
func Close() {
	_ = redisDb.Close()
}
