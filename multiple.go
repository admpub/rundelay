package rundelay

import (
	"sync"
	"time"
)

func NewMultiple[T any](delay time.Duration, f func(T) error) *Multiple[T] {
	return &Multiple[T]{
		mp:    map[string]RunDelayer[T]{},
		mu:    sync.RWMutex{},
		delay: delay,
		exec:  f,
	}
}

type Multiple[T any] struct {
	mp    map[string]RunDelayer[T]
	mu    sync.RWMutex
	delay time.Duration
	exec  func(T) error
}

func (m *Multiple[T]) get(k string) (RunDelayer[T], bool) {
	m.mu.RLock()
	val, ok := m.mp[k]
	m.mu.RUnlock()
	return val, ok
}

func (m *Multiple[T]) set(k string, v RunDelayer[T]) {
	m.mu.Lock()
	m.mp[k] = v
	m.mu.Unlock()
}

func (m *Multiple[T]) Delete(k string) {
	if m.mu.TryLock() {
		delete(m.mp, k)
		m.mu.Unlock()
		return
	}

	delete(m.mp, k)
}

func (m *Multiple[T]) Init(delay time.Duration, f func(T) error) {
	m.exec = f
	m.delay = delay
}

func (m *Multiple[T]) Run(k string, v T) bool {
	m.mu.Lock()
	val, ok := m.mp[k]
	if !ok {
		val = New(m.delay, func(t T) (err error) {
			err = m.exec(t)
			//m.Delete(k)
			return
		})
		m.mp[k] = val
	}
	m.mu.Unlock()
	return val.Run(v)
}

func (m *Multiple[T]) Done(k string) error {
	val, ok := m.get(k)
	if ok {
		return val.Done()
	}
	return nil
}

func (m *Multiple[T]) Close() (err error) {
	m.mu.Lock()
	for k, v := range m.mp {
		err = v.Close()
		if err != nil {
			break
		}
		delete(m.mp, k)
	}
	m.mu.Unlock()
	return
}

func (m *Multiple[T]) Range(cb func(string, RunDelayer[T])) error {
	m.mu.RLock()
	for k, v := range m.mp {
		cb(k, v)
	}
	m.mu.RUnlock()
	return nil
}
