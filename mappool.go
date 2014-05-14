package ujson

import "sync/atomic"

type MapPool struct {
	misses int32
	list   chan *MapItem
}

type MapItem struct {
	pool *MapPool
	m    map[string]interface{}
}

func newMapItem(p *MapPool) *MapItem {
	return &MapItem{
		pool: p,
		m:    make(map[string]interface{}),
	}
}

func (item *MapItem) Close() error {
	for k := range item.m {
		delete(item.m, k)
	}
	if item.pool != nil {
		item.pool.list <- item
	}
	return nil
}

func NewMapPool(count int) *MapPool {

	p := &MapPool{
		list: make(chan *MapItem, count),
	}
	for i := 0; i < count; i++ {
		p.list <- newMapItem(p)
	}
	return p
}

func (pool *MapPool) Checkout() *MapItem {
	var m *MapItem
	select {
	case m = <-pool.list:
	default:
		atomic.AddInt32(&pool.misses, 1)
		m = newMapItem(pool)
	}
	return m
}

func (pool *MapPool) Len() int {
	return len(pool.list)
}

func (pool *MapPool) Misses() int32 {
	return atomic.LoadInt32(&pool.misses)
}
