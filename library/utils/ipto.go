package utils

import (
	"fmt"
	"github.com/ipipdotnet/ipdb-go"
)

// IpToCountry 根据IP获取国家名
func IpToCountry(ipStr string) string {
	if ipStr == "" {
		return "未知"
	}
	appPath, _ := GetAppPath()
	db, err := ipdb.NewCity(fmt.Sprintf("%s/asserts/ipipfreedb/ipipfree.ipdb", appPath))
	//db, err := ipdb.NewCity("./ipipfreedb/ipipfree.ipdb")
	if err != nil {
		return "未知"
	}
	// 更新 ipdb 文件后可调用 Reload 方法重新加载内容
	//if flag {
	//	db.Reload("./ipipfreedb/ipipfree.ipdb")
	//}
	//查找IP信息
	var IpInfo map[string]string
	IpInfo, err = db.FindMap(ipStr, "CN")
	if err != nil {
		return "未知"
	}
	return IpInfo["country_name"]
}
