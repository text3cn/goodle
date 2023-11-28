package timekit

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/providers/goodlog"
	"regexp"
	"strconv"
	"strings"
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

// 格式化时间戳
// Date("2006-01-02 15:04:05", 1524799394)
// Date("2006/01/02 15:04:05 PM", 1524799394)
func Date(format string, timestamp int64) string {
	return time.Unix(timestamp, 0).Format(format)
}

// 2006-01-02T15:04:05+08:00 转 2006-01-02 15:04:05
func DateTimeFormat(timetime interface{}) string {
	var ret string
	if time, ok := timetime.(time.Time); ok {
		ret = time.Format("2006-01-02 15:04:05")
	}
	return ret
}

// 验证是否是一个格里高里日期
func Checkdate(month, day, year int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		// leap year
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}
	return true
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

// 字符串时间转时间戳
// 例如：Strtotime("2006-01-02 15:04:05", strtime)
// 例如：Strtotime("2006-01-02", strtime)
func Strtotime(format, strtime string) (int64, error) {
	t, err := time.ParseInLocation(format, strtime, time.Local)
	if err != nil {
		return 0, err
	}
	timezone, _ := time.LoadLocation(TimeZone)
	return t.In(timezone).Unix(), nil
}

func TimestampToDatetimeStr(timestamp int64) string {
	timeobj := time.Unix(int64(timestamp), 0)
	date := timeobj.Format("2006-01-02 15:04:05")
	return date
}

// 字符串时间转字符串日期
func DatetimeStrToDateStr(datetimeStr string) string {
	timestamp, _ := Strtotime("2006-01-02 15:04:05", datetimeStr)
	return time.Unix(timestamp, 0).Format("2006-01-02")
}

// 将 mysql 的 datetime 类型的时间字符串转为时间戳
func Datetime2Timestamp(datetime string) int {
	// 使用parseInLocation将字符串格式化返回本地时区时间, 同 php 的 strtotime()
	stamp, _ := time.ParseInLocation("2006-01-02 15:04:05", datetime, time.Local)
	return int(stamp.Unix())
}

// 将中间带 T 的时间字符串转为 time.Time
func DatetimeT2Time(datetime string) time.Time {
	time, _ := time.ParseInLocation(time.RFC3339, datetime, time.Local)
	return time
}

// 将时间戳转成 mysql datetime
func TimestampToDatetime(tiemstamp int64) string {
	return Date("2006-01-02 15:04:05", tiemstamp)
}

// 时间戳转日期字符串
func TimestampToDate(tiemstamp int64) string {
	return Date("2006-01-02", tiemstamp)
}

// 2019-01-01 15:22:22 格式字符串转 time.Time
func Str2time(strtime string) (parsedTime time.Time) {
	// 定义日期时间格式
	dateTimeFormat := "2006-01-02 15:04:05"
	// 调整时差
	local := time.FixedZone("Local", FixedZone)
	// 将字符串解析为 time.Time 类型
	var err error
	parsedTime, err = time.ParseInLocation(dateTimeFormat, strtime, local)
	if err != nil {
		goodlog.Error("Error:", err)
	}
	return
}

// 传入字符串日期，加或减去 n 天
// date string 日期格式：2202-03-01
// date int 传正数为加，负数为减去 n 天
func DateStrAddDay(date string, n int) string {
	timestamp, err := Strtotime("2006-01-02", date)
	if err != nil {
		fmt.Println(err)
	}
	timestamp += int64((3600 * 24) * n)
	return TimestampToDate(timestamp)
}

// 将 2023-11-27T21:10:10+07:00 转 2023-11-27 21:10:10 形式
// 直接通过字符串截取形式完成，不管时区
func Time2Str(t time.Time) string {
	timeString := cast.ToString(t)
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	matches := re.FindStringSubmatch(timeString)
	if len(matches) > 1 {
		return matches[1]
	}
	return timeString
}

// 将 2023-11-27T21:10:10+07:00 转 2023-11-27 21:10:10 形式
// 直接通过字符串截取形式完成，不管时区
func TimeStr2Str(timeString string) string {
	timeString = strings.Replace(timeString, "T", " ", 1)
	seg := strings.Split(timeString, "+")
	return seg[0]
}
