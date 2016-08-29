package net

import (
	"atman/mm"
	"atman/xen"
)

type Buffer struct {
	ID   int
	Gref xen.Gref
	*mm.Page
}

type BufferPool struct {
	buffers []Buffer
	free    []int
}

func NewBufferPool(size int) *BufferPool {
	pool := &BufferPool{
		buffers: make([]Buffer, size),
		free:    make([]int, size),
	}

	for i := 0; i < size; i++ {
		pool.buffers[i].ID = i
		pool.free[i] = i
	}

	return pool
}

func (p *BufferPool) Lookup(id int) *Buffer {
	return &p.buffers[id]
}

func (p *BufferPool) Get() (*Buffer, bool) {
	id, ok := p.getID()
	if !ok {
		return nil, false
	}

	buffer := &p.buffers[id]
	if buffer.Page == nil {
		buffer.Page = mm.AllocPage()
	}
	return buffer, true
}

func (p *BufferPool) getID() (int, bool) {
	if len(p.free) == 0 {
		return 0, false
	}

	id := p.free[len(p.free)-1]
	p.free = p.free[:len(p.free)-1]
	return id, true
}

func (p *BufferPool) Put(b Buffer) {
	p.free = append(p.free, b.ID)
}
