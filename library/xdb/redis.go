package xdb

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

var redisLock sync.Mutex
var rdsMap map[string]*redis.Pool = nil

func connRds(mod string) *redis.Pool {
	db := &redis.Pool{
		MaxIdle:     getConfInt(mod, "maxidle", 8),
		IdleTimeout: time.Second * time.Duration(getConfInt(mod, "maxlife", 600)),
		MaxActive:   getConfInt(mod, "maxopen", 16),
		Wait:        true,

		Dial: func() (redis.Conn, error) {
			host := getConfString(mod, "host")
			passwd := getConfString(mod, "passwd")
			logs.Info("redis connect host=" + host)
			c, err := redis.Dial("tcp", host)
			if nil != err {
				logs.Error("connect redis error host:%s msg:%s", host, err.Error())
				return nil, err
			}

			if "" != passwd {
				if _, err := c.Do("AUTH", passwd); nil != err {
					logs.Error("connect redis auth error host:%s passwd:%s msg:%s", host, passwd, err.Error())
					if err := c.Close(); nil != err {
						logs.Error("close redis %s", err.Error())
					}
					return nil, err
				}
			}
			return c, nil
		},
	}
	return db
}

func Rds(mod string) *redis.Pool {
	if nil == rdsMap {
		logs.Error("xdb.Rds rdsMap is nil")
		return nil
	}
	redisLock.Lock()
	conn := rdsMap[mod]
	if nil == conn {
		conn = connRds(mod)
		rdsMap[mod] = conn
	}
	redisLock.Unlock()
	return conn
}
