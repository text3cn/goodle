package timekit

import (
	"github.com/spf13/cast"
	"github.com/text3cn/goodle/providers/goodlog"
	"regexp"
	"strings"
	"time"
)

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

func Time2stamp(t time.Time) int64 {
	unixTimestampMillis := t.UnixNano() / int64(time.Millisecond)
	return unixTimestampMillis
}
