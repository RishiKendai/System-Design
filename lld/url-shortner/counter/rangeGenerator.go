package counter

import "sync"

type RangeGenerator struct {
	ID    string
	mu    sync.Mutex
	next  int64 // next value to hand out (inclusive)
	limit int64 // exclusive upper bound of current range
	block int64
}

func NewRangeGenerator(id string, block int64) *RangeGenerator {
	if block < 1 {
		block = 1
	}
	start := getAllocator().Allocate(block)
	return &RangeGenerator{
		ID:    id,
		next:  start,
		limit: start + block,
		block: block,
	}
}

func (g *RangeGenerator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.next >= g.limit {
		start := getAllocator().Allocate(g.block)
		g.next = start
		g.limit = start + g.block
	}
	v := g.next
	g.next++
	return v
}
