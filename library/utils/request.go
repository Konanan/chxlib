package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//IsGet 判断是否为Get请求
func IsGet(r *http.Request) bool {
	return r.Method == "GET"
}

//Response 返回json
func Response(w http.ResponseWriter, msg interface{}, code int) {
	dist := map[string]interface{}{
		"msg":  msg,
		"code": code,
	}

	b, err := json.Marshal(&dist)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	w.Header().Set("Content-Type: ", "application/json;charset=utf-8")
	_, _ = w.Write([]byte(b))
}

//layui 后台返回需要的json格式
func LayuiJson(code int, msg string, data, count interface{}) (jsonData map[string]interface{}) {
	jsonData = make(map[string]interface{}, 3)
	jsonData["code"] = code
	if msg == "" {
		jsonData["data"] = data
		jsonData["count"] = count
	} else {
		jsonData["msg"] = msg
	}
	jsonData["time_stamp"] = time.Now()
	return
}
