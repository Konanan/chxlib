package xdb

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/elastic/go-elasticsearch/v7"
)

var esLock sync.Mutex
var esMap map[string]*elasticsearch.Client = nil

func connEs(mod string) *elasticsearch.Client {
	logs.Info(getConfString(mod, "host"), getConfInt(mod, "maxidle", 8), getConfInt(mod, "maxidle", 8), getConfString(mod, "uname"), getConfString(mod, "passwd"))
	conf := elasticsearch.Config{
		Addresses: []string{
			getConfString(mod, "host"),
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   getConfInt(mod, "maxidle", 8),
			MaxIdleConns:          getConfInt(mod, "maxidle", 8),
			IdleConnTimeout:       time.Second * time.Duration(getConfInt(mod, "maxlife", 600)),
			ResponseHeaderTimeout: time.Second * 10,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MaxVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
		Username: getConfString(mod, "uname"),
		Password: getConfString(mod, "passwd"),
	}
	es, err := elasticsearch.NewClient(conf)
	if nil != err {
		logs.Error("connect es error msg:%s", err.Error())
		return nil
	}
	return es
}

func Es(mod string) *elasticsearch.Client {
	if nil == esMap {
		logs.Error("xdb.Db esMap is nil")
		return nil
	}
	esLock.Lock()
	conn := esMap[mod]
	if nil == conn {
		conn = connEs(mod)
		esMap[mod] = conn
	}
	esLock.Unlock()
	return conn
}
