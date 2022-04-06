package redis_v1

import (
	"github.com/gomodule/redigo/redis"
	"github.com/larspensjo/config"
	log "github.com/sirupsen/logrus"
	"time"
)

// 结构体1：
type redisRun struct {
	dbInstanceName string      // 名称:default等
	handle         *redis.Pool // redis连接池句柄
}

// 全局变量
var (
	handles    = make(map[string]*redisRun) // 实例句柄
	pathConfig = ""                         // 配置文件地址
)

// 初始化函数
func Init(cfgPath string) {
	pathConfig = cfgPath
}

// 获取连接池句柄
func Handle(dbInstance ...string) *redis.Pool {
	// 没有参数就用 default
	instance := "default"
	if len(dbInstance) > 0 {
		instance = dbInstance[0]
	}
	
	// 不存在创建
	if _, ok := handles[instance]; !ok {
		handles[instance] = new(redisRun)
		handles[instance].dbInstanceName = instance
		handles[instance].handle = newRedisPool(instance)
		return handles[instance].handle
	} else {
		// 存在返回
		return handles[instance].handle
	}
}

// 创建redis连接池
func newRedisPool(dbInstance string) *redis.Pool {
	var (
		cfg *config.Config
		
		address     string
		pass        string
		db          int
		MaxIdle     int
		MaxActive   int
		IdleTimeout int
		
		err error
	)
	// 读取配置文件
	cfg, err = config.ReadDefault(pathConfig)
	if err != nil {
		log.Error("读取配置文件出错", err)
		return nil
	}
	address, _ = cfg.String("redis_"+dbInstance, "address") // 127.0.0.1:6379
	pass, _ = cfg.String("redis_"+dbInstance, "pass")       // nlsper
	db, _ = cfg.Int("redis_"+dbInstance, "db")
	MaxIdle, _ = cfg.Int("redis_"+dbInstance, "MaxIdle")
	MaxActive, _ = cfg.Int("redis_"+dbInstance, "MaxActive")
	IdleTimeout, _ = cfg.Int("redis_"+dbInstance, "IdleTimeout")
	
	return &redis.Pool{
		MaxIdle:     MaxIdle,
		MaxActive:   MaxActive,
		IdleTimeout: time.Second * time.Duration(IdleTimeout),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp",
				address,
				redis.DialDatabase(db),
				redis.DialPassword(pass),
			)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			// 2、访问认证
			// if _, err =c.Do("AUTH",redisPass);err!=nil{
			//	c.Close()
			//	return nil,err
			// }
			return c, nil
		},
		// 定时检查redis是否出状况
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}
