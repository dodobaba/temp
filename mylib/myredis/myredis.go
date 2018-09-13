package myredis

import (
	"shoppingzone/conf"
	"shoppingzone/mylib/mylog"
	"time"

	"github.com/gomodule/redigo/redis"
)

const redisURL = "127.0.0.1"

var redisSession *redis.Pool

//GetRedisSession :
func GetRedisSession() redis.Conn {
	if redisSession == nil {
		redisSession = redis.NewPool(
			func() (redis.Conn, error) {
				return redis.Dial(conf.RedisConf.Protocol, conf.RedisConf.Host)
			}, 3)
		redisSession.IdleTimeout = time.Duration(time.Second * 240)
		redisSession.MaxActive = conf.RedisConf.Poolconf.MaxActive
		redisSession.MaxIdle = conf.RedisConf.Poolconf.MaxIdle
	}
	return redisSession.Get()
}

//Set :
func Set(key string, value string, expire int) error {
	conn := GetRedisSession()
	defer conn.Close()
	conn.Send("SET", key, value)
	conn.Send("EXPIRE", key, expire)
	err := conn.Flush()
	if err != nil {
		mylog.Tf("[Error]", "MyRedis", "Set", "cache not saved %s", err.Error())
		return err
	}
	return nil
}

//Get :
func Get(key string) (string, error) {
	conn := GetRedisSession()
	defer conn.Close()
	r, err := redis.String(conn.Do("GET", key))
	if err != nil {
		mylog.Tf("[Error]", "MyRedis", "Get", "%s", err.Error())
		return "", err
	}
	return r, nil
}
