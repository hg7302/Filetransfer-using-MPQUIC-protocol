package utils

import "sync/atomic"

// An AtomicBool is an atomic bool
type AtomicBool struct {
	v int64
}

// Set sets the value
func (a *AtomicBool) Set(value bool) {
	var n int64
	if value {
		n = 1
	}
	atomic.StoreInt64(&a.v, n)
}

// Get gets the value
func (a *AtomicBool) Get() bool {
	return atomic.LoadInt64(&a.v) != 0
}
