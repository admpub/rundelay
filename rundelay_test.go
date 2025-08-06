package rundelay

import (
	"math"
	"sync"
	"testing"
	"time"
)

func _TestRunDelay(t *testing.T) {
	target := time.Now()
	var execTime time.Time
	exec := func(i int) error {
		execTime = time.Now()
		t.Logf(`%d.========= test delay run: %v`, i, execTime.Format(time.RFC3339Nano))
		return nil
	}
	delay := time.Second * 2
	dr := New(delay, exec)
	dr.recvNotify = func(i int, tgt time.Time) {
		t.Logf(`%d. recv notify, new target: %v`, i, tgt.Format(time.RFC3339Nano))
		target = time.Now()
	}
	wg := sync.WaitGroup{}

	t.Logf(`test delay start: %v`, target.Format(time.RFC3339Nano))

	testf := func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				if i%2 == 0 {
					time.Sleep(time.Millisecond * 1000)
				} else {
					time.Sleep(time.Millisecond * 500)
				}
				defer wg.Done()
				ok := dr.Run(i)
				if !ok {
					t.Logf(`~~~~~ skipped %d`, i)
				}
			}(i)
		}
	}

	testf()

	time.Sleep(time.Second * 3)

	testf()

	wg.Wait()
	err := dr.Done()
	if err != nil {
		panic(err)
	}
	vdur := execTime.Sub(target)
	if int64(math.Round(vdur.Seconds())) != int64(math.Round(delay.Seconds())) {
		panic(execTime.Format(time.RFC3339Nano) + ` - ` + target.Format(time.RFC3339Nano) + ` = ` + vdur.String() + ` != ` + delay.String())
	}

	time.Sleep(time.Second * 3)

	testf()

	wg.Wait()

	err = dr.Done()
	if err != nil {
		panic(err)
	}
	vdur = execTime.Sub(target)
	if int64(math.Round(vdur.Seconds())) != int64(math.Round(delay.Seconds())) {
		panic(execTime.Format(time.RFC3339Nano) + ` - ` + target.Format(time.RFC3339Nano) + ` = ` + vdur.String() + ` != ` + delay.String())
	}

	dr.Close()
}
