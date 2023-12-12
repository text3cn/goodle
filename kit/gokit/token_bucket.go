package gokit

import (
	"time"
)

type restrictor struct {
	tokenBucket  chan struct{}
	fillInterval time.Duration
	cap          int
}

// fillInterval 每隔多久向桶中填充一个令牌，例如：time.Microsecond * 10 每 10 毫秒填充一个令牌
// cap 令牌桶容量，例如: 100。两个参数合起来，每 10 毫秒填充一个令牌，1 秒刚好填满 100 个，也就是 100 的 qps
func NewTokenBucket(fillInterval time.Duration, cap int) restrictor {
	s := restrictor{
		fillInterval: fillInterval,
		cap:          cap,
	}
	s.fillToken()
	return s
}

// 填充令牌
func (this *restrictor) fillToken() {
	tokenBucket := make(chan struct{}, this.cap)
	go func() {
		// 这里可以看到 Go 的定时器存在大约 0.001 秒的误差
		// 所以如果令牌桶大小在 1000 以上的填充可能会有一定的误差。对一般的服务来说，这一点误差无关紧要。
		ticker := time.NewTicker(this.fillInterval)
		for {
			select {
			case <-ticker.C:
				select {
				// 填充令牌
				case tokenBucket <- struct{}{}:
					//fmt.Println("填充一个令牌，总共: ", len(tokenBucket), " 个")
				default:
				}
				// fmt.Println("current token count: ", len(tokenBucket), time.Now())
			}
		}
	}()
	this.tokenBucket = tokenBucket
}

// 取令牌
// await 是否等待令牌，默认 false
// await 为 true 时如果桶中没有令牌则会阻塞线程，等待桶中有令牌的时候才返回 true，即把流量挡着等有空位一个放一个进来
// await 为 false 时直接返回现在有没有令牌，例如没有令牌时调用者可以直接给用户返回"请稍后再试"的提示
func (this *restrictor) TakeToken(await ...bool) bool {
	block := false
	if len(await) > 0 {
		block = await[0]
	}
	var takeResult bool
	if block {
		select {
		case <-this.tokenBucket:
			takeResult = true
		}
	} else {
		select {
		case <-this.tokenBucket:
			// fmt.Println("取出一个令牌，还剩: ", len(this.tokenBucket), " 个")
			takeResult = true
		default:
			// fmt.Println("没有可用令牌")
			takeResult = false
		}
	}
	return takeResult
}
