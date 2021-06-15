package atomics

import "sync/atomic"

type AtomicBool struct {
	value int32
}

func (ab *AtomicBool) Get() bool {
	return atomic.LoadInt32(&ab.value)&1 != 0
}

func (ab *AtomicBool) Set(value bool) {
	if value {
		atomic.StoreInt32(&ab.value, 1)
	} else {
		atomic.StoreInt32(&ab.value, 0)
	}
}
