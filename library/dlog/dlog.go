package chxlib

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func Init() {
	process, err := getProcessName()
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	path := beego.AppConfig.DefaultString("log::path", "")
	depth := beego.AppConfig.DefaultInt("log::depth", 3)
	rotate := beego.AppConfig.DefaultString("log::rotate", "true")
	daily := beego.AppConfig.DefaultString("log::daily", "true")
	hourly := beego.AppConfig.DefaultString("log::hourly", "false")
	level := beego.AppConfig.DefaultInt("log::level", 6)

	//拼接日志文件目录
	file := fmt.Sprintf("%s%s.log", path, process)

	logs.SetPrefix(fmt.Sprintf("[%d]", os.Getpid()))
	logs.SetLogFuncCallDepth(depth)
	logs.SetLevel(level)

	logConfig := fmt.Sprintf(`{"filename":"%s","rotate":%s,"maxdays":30,"maxsize":0,"maxlines":0,"daily":%s,"hourly":%s,"rotateperm":"0644","perm":"0644"}`, file, rotate, daily, hourly)
	if err := logs.SetLogger(logs.AdapterFile, logConfig); nil != err {
		fmt.Println("init log err ", err.Error())
		return
	}
	logs.Info("init log ok : %s", logConfig)
}

func getProcessName() (string, error) {
	args := os.Args[0]
	if len(args) <= 0 {
		return "", errors.New("InitWithProcess os.Args <= 0")
	}
	names := strings.Split(args, "/")
	if len(names) <= 0 {
		return "", errors.New("InitWithProcess strings.Split <= 0")
	}
	return names[len(names)-1], nil
}
