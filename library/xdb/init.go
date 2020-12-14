package xdb

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
)

func init() {
	dbMap = make(map[string]*gorm.DB)
	rdsMap = make(map[string]*redis.Pool)
	esMap = make(map[string]*elasticsearch.Client)
}

func getConfString(name, key string) string {
	return beego.AppConfig.String(fmt.Sprintf("%s::%s", name, key))
}

func getConfInt(name, key string, def int) int {
	return beego.AppConfig.DefaultInt(fmt.Sprintf("%s::%s", name, key), def)
}
