package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/loudbund/go-redis/redis_v1"
	log "github.com/sirupsen/logrus"
	"time"
)

func init() {
	log.SetReportCaller(true)
	redis_v1.Init("test.conf")
}

func main() {
	// redis连接池获取
	conn := redis_v1.Handle().Get()
	defer func() { _ = conn.Close() }()

	// 1、设置键值内容 SET
	if d, err := redis.String(conn.Do("SET", "wawa", "娃哈")); err != nil {
		fmt.Println(d, err)
	}

	// 2、获取键值内容 GET
	if d, err := redis.String(conn.Do("GET", "wawa")); err != nil {
		fmt.Println(1, err)
	} else {
		// 返回内容
		fmt.Println(2, d)
	}

	// 3、设置键值带过期时间，过期时间单位(秒) SET
	if d, err := redis.String(conn.Do("SET", "haha", "哈哈", "EX", 60)); err != nil {
		fmt.Println(3, err)
	} else {
		// 返回OK
		fmt.Println(4, d)
	}

	// 4、写入队列, RPUSH
	if d, err := conn.Do("RPUSH", "zList", "x", "y"); err != nil {
		fmt.Println(5, err)
	} else {
		// 返回写成功数量
		fmt.Println(6, d)
	}

	// 5、读取队列 BLPOP
	go func() {
		for {
			if value, err := conn.Do("BLPOP", "zList", 1); err != nil {
				fmt.Println(7, err)
			} else {
				if value != nil {
					fmt.Println(8, string(value.([]interface{})[0].([]byte)))
					fmt.Println(8, string(value.([]interface{})[1].([]byte)))
				} else {

				}
			}
		}
	}()

	// 	延时5秒退出
	time.Sleep(time.Second * 5)
}
