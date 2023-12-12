package timekit

import (
	"github.com/spf13/cast"
	"strconv"
	"time"
)

var TimeZone = "Asia/Shanghai"

// 转 time.Time 类型时会自动使用 UTC 时间，比北京时间晚 8 小时，
// FixedZone 基于服务器所在时区进行时区偏移
var FixedZone = 3600 * 8

// 获取指定时区的当前时间字符串，可直接存入 mysql datetime
// args[0] 时区，如：Asia/Ho_Chi_Minh
func NowDatetimeStr(args ...string) string {
	var timezone *time.Location
	if len(args) > 0 {
		TimeZone = args[0]
	}
	timezone, _ = time.LoadLocation(TimeZone)
	t := time.Now().In(timezone).Format("2006-01-02 15:04:05")
	return cast.ToString(t)
}

func DatetimeStr() string {
	timezone, _ := time.LoadLocation(TimeZone)
	int64 := time.Now().In(timezone).Unix()
	timestamp := time.Unix(int64, 0).Format("2006-01-02 15:04:05")
	return timestamp
}

// args[0] 往前或往后推几天，昨天传 -1
func DateStr(args ...int64) string {
	var timestamp int64
	timezone, _ := time.LoadLocation(TimeZone)
	timestamp = time.Now().In(timezone).Unix()
	if len(args) > 0 {
		timestamp += args[0] * (3600 * 24)
	}
	str := time.Unix(timestamp, 0).Format("2006-01-02")
	return str
}

// 获取年月，格式：2202（2022年2月）
func YearMonthShortStr() string {
	timezone, _ := time.LoadLocation(TimeZone)
	int64 := time.Now().In(timezone).Unix()
	timestamp := time.Unix(int64, 0).Format("0601")
	return timestamp
}

// 获取当前时间的 mysq datetime 格式 （ 2006-01-02 15:04:05 ）
func Datetime() string {
	timezone, _ := time.LoadLocation(TimeZone)
	int64 := time.Now().In(timezone).Unix()
	timestamp := time.Unix(int64, 0).Format("2006-01-02 15:04:05")
	return timestamp
}

// 获取今天零点的时间戳
func TodayStartTime() int64 {
	timeStamp := time.Now()
	newTime := time.Date(timeStamp.Year(), timeStamp.Month(), timeStamp.Day(), 0, 0, 0, 0, timeStamp.Location())
	return newTime.Unix()
}

// 获取今天最后时间戳
func TodayEndTime() int64 {
	timeStamp := time.Now()
	newTime := time.Date(timeStamp.Year(), timeStamp.Month(), timeStamp.Day(), 23, 59, 59, 0, timeStamp.Location())
	return newTime.Unix()
}

// 获取微秒时间戳
func Microtime(args ...string) int {
	//return int(time.Now().UnixNano())
	timezone, _ := time.LoadLocation(TimeZone)
	tz := timezone
	if len(args) > 0 {
		tz, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	}
	return int(time.Now().In(tz).UnixNano())
}

// 获取当前时间戳
func NowTimestamp() int {
	return int(time.Now().Unix())
}

// 获取毫秒时间戳
// args[0] 时区
func Millisecond() int {
	return int(time.Now().UnixNano() / 1e6)
}

// 获取今天零点的时间戳
func TodayTime() int {
	timezone, _ := time.LoadLocation(TimeZone)
	timeStamp := time.Now().In(timezone)
	newTime := time.Date(timeStamp.Year(), timeStamp.Month(), timeStamp.Day(), 0, 0, 0, 0, timeStamp.Location())
	return int(newTime.Unix())
}

// 获取今天零点的时间字符串，格式：2021-11-03 00:00:00
func DateTodayZeroStr() string {
	return DateTodayStr() + " 00:00:00"
}

// 获取今天 23:59 的时间字符串，格式：2021-11-03 23:59:59
func DateToday2359Str() string {
	return DateTodayStr() + " 23:59:59"
}

// 获取当前年月日
func DateTodayInt() (int, int, int) {
	year := time.Now().Year()
	monthStr := time.Now().Format("01")
	month, _ := strconv.Atoi(monthStr)
	day := time.Now().Day()
	return year, month, day
}

// 获取今天日期字符串，格式 2021-11-03
func DateTodayStr() string {
	timezone, _ := time.LoadLocation(TimeZone)
	int64 := time.Now().In(timezone).Unix()
	str := time.Unix(int64, 0).Format("2006-01-02")
	return str
}

// 今天日期字符串，格式 220105
func DateTodayShortStr() string {
	timezone, _ := time.LoadLocation(TimeZone)
	int64 := time.Now().In(timezone).Unix()
	str := time.Unix(int64, 0).Format("060102")
	return str
}

// 获取几天前的日期字符串
// day 几天前
func DateBeforDaysStr(day int) string {
	stamp := int64(TodayTime() - day*3600*24)
	date := TimestampToDate(stamp)
	return date
}

// 获取当前时间戳 - 字符串类型
// addtime 增加时间，秒为单位。常用于返回到期时间
func TimeStampString(addtime ...int) string {
	var add int
	if len(addtime) > 0 {
		add = addtime[0]
	}
	timezone, _ := time.LoadLocation(TimeZone)
	return strconv.FormatInt(time.Now().In(timezone).Unix()+int64(add), 10)
}

// 获取本周零点时间戳
func WeekTime() int {
	var week = map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}
	timezone, _ := time.LoadLocation(TimeZone)
	timeStamp := time.Now().In(timezone)
	weekStr := timeStamp.Weekday().String()
	weekInt := week[weekStr]
	newTime := time.Date(timeStamp.Year(), timeStamp.Month(), timeStamp.Day()-weekInt, 0, 0, 0, 0, timeStamp.Location())
	return int(newTime.Unix())
}
