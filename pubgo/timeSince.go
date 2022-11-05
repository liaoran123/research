package pubgo

import (
	"fmt"
	"time"
)

//执行时间

type ts struct {
	stime time.Time
}

func Newts() ts {
	return ts{stime: time.Now()}
}
func (t *ts) Set() { //重置开始时间
	t.stime = time.Now()
}
func (t *ts) Gts() time.Duration {
	return time.Since(t.stime)
}
func (t *ts) Gstrts() string {
	return fmt.Sprintf("%v", time.Since(t.stime))
}

//fmt.Sprintf("%v", ys)
