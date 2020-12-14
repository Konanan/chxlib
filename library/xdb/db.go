package xdb

import (
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbLock sync.Mutex
var dbMap map[string]*gorm.DB = nil

func Db(mod string) *gorm.DB {
	if nil == dbMap {
		logs.Error("xdb.Db dbMap is nil")
		return nil
	}
	dbLock.Lock()
	conn := dbMap[mod]
	if nil == conn {
		conn = connDb(mod)
		dbMap[mod] = conn
	}
	dbLock.Unlock()

	return conn
}

func connDb(mod string) *gorm.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		getConfString(mod, "uname"),
		getConfString(mod, "passwd"),
		getConfString(mod, "host"),
		getConfString(mod, "name"))

	maxIdle := getConfInt(mod, "maxidle", 8)
	maxOpen := getConfInt(mod, "maxopen", 128)
	maxLife := getConfInt(mod, "maxlife", 600)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dataSource,
	}), &gorm.Config{})
	if nil != err {
		logs.Error("open mysql fail,err:%s", err.Error())
		return nil
	}
	//Debug模式
	if beego.DEV == beego.AppConfig.String("runmode") {
		db = db.Debug()
	}
	sqlDb, err := db.DB()
	if err != nil {
		logs.Error("connect mysql fail, err:%s", err.Error())
		panic(err)
	}
	//最大连接生命周期
	sqlDb.SetConnMaxLifetime(time.Duration(maxLife) * time.Second)
	//最大连接数
	sqlDb.SetMaxOpenConns(maxOpen)
	//最大空闲连接数
	sqlDb.SetMaxIdleConns(maxIdle)

	if err := sqlDb.Ping(); nil != err {
		logs.Error("ping mysql fail, err:%s", err.Error())
		return nil
	}
	logs.Info("connect mysql success, mod:%s dataSource:%s maxIdle:%d maxOpen:%d", mod, dataSource, maxIdle, maxOpen)
	return db
}
