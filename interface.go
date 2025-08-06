package rundelay

import (
	"time"
)

var _ Rundelayer[any] = (*RunDelay[any])(nil)

type Rundelayer[T any] interface {
	Init(delay time.Duration, f func(T) error)
	Run(T) bool
	Done() error // 阻塞获取执行结果
	Close() error
}
