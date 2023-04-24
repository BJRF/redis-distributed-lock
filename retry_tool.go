package redis_distributed_lock

import "time"

// RetryTool 结构体
type RetryTool struct {
	// 重试间隔
	Interval time.Duration
	// 最大次数
	Max int
	cnt int
}

// RetryStrategy 重试策略接口，内含next方法，返回是否重置
type RetryStrategy interface {
	Next() (time.Duration, bool)
}

// Next 返回下一次重试的间隔，并判断是否还需要继续重试
func (f *RetryTool) Next() (time.Duration, bool) {
	f.cnt++
	return f.Interval, f.cnt <= f.Max
}
