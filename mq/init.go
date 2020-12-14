package chxlib

import (
	"fmt"

	"github.com/astaxie/beego"
)

func init() {
	kafkaMap = make(map[string]*Kafka)
}

func getConfString(name, key string) string {
	return beego.AppConfig.String(fmt.Sprintf("%s::%s", name, key))
}

func getConfInt(name, key string, def int) int {
	return beego.AppConfig.DefaultInt(fmt.Sprintf("%s::%s", name, key), def)
}

func getConfBool(name, key string) bool {
	return beego.AppConfig.DefaultBool(fmt.Sprintf("%s::%s", name, key), false)
}
