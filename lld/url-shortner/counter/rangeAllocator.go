package counter

import (
	"math"
	"sync"
)

var (
	once      sync.Once
	allocator *RangeAllocator
)

type RangeAllocator struct {
	mu      sync.Mutex
	counter int64
}

func NewRangeAllocator() *RangeAllocator {
	return &RangeAllocator{
		counter: int64(math.Pow(62, 6)), // ≈ 56,800,235,584
	}
}

// Allocate reserves a half-open interval [start, end): values start, start+1, …, end-1.
// count must be positive; otherwise it is treated as 1.
func (a *RangeAllocator) Allocate(block int64) int64 {
	if block < 1 {
		block = 1
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	start := a.counter
	a.counter += block
	return start
}

func getAllocator() *RangeAllocator {
	once.Do(func() {
		allocator = NewRangeAllocator()
	})
	return allocator
}
