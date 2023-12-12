package timekit

import "fmt"

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
