package common

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (once *Once) Do(f func() error) error {
	if atomic.LoadUint32(&once.done) > 0 {
		return nil
	}

	once.m.Lock()
	defer once.m.Unlock()

	if once.done > 0 {
		return nil
	}

	err := f()
	if err != nil {
		return err
	}

	atomic.StoreUint32(&once.done, 1)
	return nil
}
