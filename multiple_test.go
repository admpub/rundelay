package rundelay

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMultiple(t *testing.T) {
	target := time.Now()
	delay := time.Second * 2
	m := NewMultiple(delay, func(i int) error {
		execTime := time.Now()
		t.Logf(`%d.========= test delay run: %v`, i, execTime.Format(time.RFC3339Nano))
		return nil
	})
	defer m.Close()

	t.Logf(`test delay start: %v`, target.Format(time.RFC3339Nano))

	wg := sync.WaitGroup{}

	testf := func() {
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func(i int) {
				if i%2 == 0 {
					time.Sleep(time.Millisecond * 1000)
				} else {
					time.Sleep(time.Millisecond * 500)
				}
				defer wg.Done()
				ok := m.Run(fmt.Sprint(i%2), i)
				if !ok {
					t.Logf(`~~~~~ skipped %d`, i)
				}
			}(i)
		}
	}

	testf()

	wg.Wait()

	m.Range(func(s string, rd RunDelayer[int]) {
		err := rd.Done()
		t.Logf(`Done: %s: %v`, s, err)
	})
}
