package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
)

const timeFormat = "2006-01-02 15:04:05.999999"

func Atoi(s string) int {
	n, err := strconv.Atoi(s)
	if nil != err {
		logs.Error(err.Error())
		return 0
	}
	return n
}

// 字符串转时间戳 必须格式为 2006-01-02 15:04:05
func TimeStrToInt(timeStr string) int64 {
	timeInt, err := ParseDateTime(timeStr)
	if err != nil {
		return 0
	}
	return timeInt.Unix()
}

// 获取指定时间的当天0点时间
func DateStrToInt(timeStr string) int64 {
	loc, _ := time.LoadLocation("Local")
	timeInt, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
	return timeInt.Unix()
}

// 时间戳转日期
func TimeIntToStr(timestamp int) string {
	timeTemplate := "2006-01-02 15:04:05" //常规类型
	return time.Unix(int64(timestamp), 0).Format(timeTemplate)
}

// 时间戳转日期(PC)
func TimeIntToDate(timestamp int) string {
	timeTemplate := "2006-01-02" //常规类型
	return time.Unix(int64(timestamp), 0).Format(timeTemplate)
}

// 时间戳转日期到小时(PC)
func TimeIntToDateHour(baseTime int) int64 {
	curTime := time.Unix(int64(baseTime), 0)
	hourTime := curTime.Unix() - int64(curTime.Second()) - int64((60 * curTime.Minute()))
	return hourTime
}

// 时间戳转小时(PC)
func TimeIntToHour(timestamp int) string {
	timeTemplate := "15" //常规类型
	return time.Unix(int64(timestamp), 0).Format(timeTemplate)
}

//获取当天的零点时间(PC)
func GetTodayTime() int64 {
	todayDateStr := time.Now().Format("2006-01-02 00:00:00")
	curTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayDateStr, time.Local)
	return curTime.Unix()
}

//获取某一时间的零点时间(PC)
func GetDayStartUnix(baseTime int64) int64 {
	dateStr := time.Unix(baseTime, 0).Format("2006-01-02 00:00:00")
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr, time.Local)
	return startTime.Unix()
}

func ParseDateTime(str string) (t time.Time, err error) {
	loc, _ := time.LoadLocation("Local")
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}

	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}

//获取当前月的第一秒的时间
func GetFirstOfMonth() time.Time {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	return firstOfMonth
}

func TimeFormat(timestamp int64, template string) string {
	return time.Unix(timestamp, 0).Format(template)
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}

//获取项目工程路径
func GetAPPRootPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	p, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	return filepath.Dir(p)
}

//获取APP的路径
func GetAppPath() (path string, err error) {
	return beego.AppConfig.String("appPath"), nil
}

//字符串截取(PC)
func Substr(str string, start uint, length int) string {
	if start < 0 || length < -1 {
		return str
	}
	switch {
	case length == -1:
		return str[start:]
	case length == 0:
		return ""
	}
	end := int(start) + length
	if end > len(str) {
		end = len(str)
	}
	return str[start:end]
}

//将一天拆分为分钟格式
func FormatDayToMinute() (timePoint []string, timeArr map[string]int32) {
	var hour, minute string
	timeArr = make(map[string]int32, 0)
	for i := 0; i < 24; i++ {
		for j := 0; j < 60; j++ {
			if i < 10 {
				hour = "0" + strconv.Itoa(i)
			} else {
				hour = strconv.Itoa(i)
			}
			if j < 10 {
				minute = "0" + strconv.Itoa(j)
			} else {
				minute = strconv.Itoa(j)
			}
			timeStr := hour + ":" + minute
			timePoint = append(timePoint, timeStr)
			timeArr[timeStr] = 0
		}
	}
	return timePoint, timeArr
}

func JsonToString(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		logs.Error("Marshal json to string faild:" + err.Error())
	}
	return string(str)
}

//获取月份倒序(PC)
func GetDescMonthByTime(startTime, endTime int64) []string {
	if startTime == 0 || endTime == 0 {
		return nil
	}

	sTime := time.Unix(startTime, 0)
	sYear, sMonth, _ := sTime.Date()
	sLocation := sTime.Location()
	startDate := time.Date(sYear, sMonth, 1, 0, 0, 0, 0, sLocation).Unix()

	eTime := time.Unix(endTime, 0)
	eYear, eMonth, _ := eTime.Date()
	eLocation := eTime.Location()
	endDate := time.Date(eYear, eMonth, 1, 0, 0, 0, 0, eLocation).Unix()

	var months []string
	for startDate <= endDate {
		month := time.Unix(endDate, 0).Format("200601")
		months = append(months, month)
		endDate = time.Unix(endDate, 0).AddDate(0, -1, 0).Unix()
	}
	return months
}

func GetFirstDateOfMonth(d time.Time) time.Time {
	var timezone = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(timezone)
}

func GetLastDateOfMonth(d time.Time) time.Time {
	var timezone = GetFirstDateOfMonth(d).AddDate(0, 1, -1)
	return GetZeroTime(timezone)
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUUID() string {
	v4 := uuid.NewV4()
	return MD5((v4).String())
}
