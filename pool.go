package comethandler

import (
	"errors"
	"sync"
	"time"
)

const (
	DefaultMaxPoolSize = 20000
	DefaultTimeout     = 5 * time.Second
)

var (
	PoolOverFlowError = errors.New("The Context Pool is Overflow")
	PoolTimeoutError  = errors.New("The Context Pool is Timeout")
)

type ContextPool struct {
	Pool    chan *Context
	MaxSize int
	Timeout time.Duration
	Length  int
	mutex   *sync.RWMutex
}

func (p *ContextPool) timeOut(d time.Duration) chan bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(d)
		timeout <- true
	}()
	return timeout
}

func (p *ContextPool) Put(context *Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.Length+1 > p.MaxSize {
		return PoolOverFlowError
	}
	timeout := p.timeOut(p.Timeout)
	select {
	case p.Pool <- context:
		p.Length++
	case <-timeout:
		return PoolTimeoutError
	}
	return nil
}

func (p *ContextPool) Get() (context *Context, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	timeout := p.timeOut(p.Timeout)
	select {
	case context = <-p.Pool:
		p.Length--
	case <-timeout:
		err = PoolTimeoutError
	}
	return
}

// NewContextPool returns a context pool.
// size : max size of the poll
func NewContextPool(params ...interface{}) *ContextPool {
	size := DefaultMaxPoolSize
	if len(params) > 0 {
		size = params[0].(int)
	}
	return &ContextPool{
		Pool:    make(chan *Context, size),
		MaxSize: size,
		Timeout: DefaultTimeout,
		mutex:   &sync.RWMutex{},
	}
}
