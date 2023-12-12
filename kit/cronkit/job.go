package fn

import "time"

// 每天执行一次的计划任务
func CronJobDaily(fn func(), hour, min, sec int) {
	for {
		// 获取当前时间，给 next 用
		now := time.Now()
		// 通过 now 偏移 24 小时
		next := now.Add(time.Hour * 24) //time.Hour * 24
		// 然后获取下一个执行时间
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, min, sec, 0, next.Location())
		// 再计算当前时间到下一个执行时间的间隔，设置一个定时器
		t := time.NewTimer(next.Sub(time.Now()))
		<-t.C
		// 执行定时任务
		fn()
	}
}

