package comethandler

import (
	"container/list"
	"sync"
)

type ContextPool struct {
	List  *list.List
	mutex *sync.Mutex
}

func (p *ContextPool) Len() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.List.Len()
}

func (p *ContextPool) Put(context *Context) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.List.PushBack(context)
}

func (p *ContextPool) Get() *Context {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	e := p.List.Front()
	p.List.Remove(e)
	context, _ := e.Value.(*Context)
	return context
}

// NewContextPool returns a context pool.
// size : max size of the poll
func NewContextPool(params ...interface{}) *ContextPool {
	return &ContextPool{
		List:  list.New(),
		mutex: &sync.Mutex{},
	}
}
