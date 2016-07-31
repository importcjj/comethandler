package comethandler

import "net/http"
import "runtime"

type CometHandler struct {
	MessageBox chan []byte
	Pool       *ContextPool
}

func (c *CometHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	context := NewContext(rw, r)
	c.Pool.Put(context)
	// if err != nil {
	// 	rw.Write([]byte("Connection reset by peer"))
	// 	return
	// }
	runtime.Gosched()
	context.Wait()
}

func (c *CometHandler) Func(rw http.ResponseWriter, r *http.Request) {
	c.ServeHTTP(rw, r)
}

func (c *CometHandler) Broadcast(body []byte) {
	c.MessageBox <- body
}

func (c *CometHandler) handle() {
	for {
		message := <-c.MessageBox
		clientsNum := c.Pool.Len()
		for clientsNum > 0 {
			context := c.Pool.Get()
			clientsNum--
			context.Write(message)
		}
	}
}

func New() *CometHandler {
	cometHandler := new(CometHandler)
	cometHandler.MessageBox = make(chan []byte, 100)
	cometHandler.Pool = NewContextPool()

	go cometHandler.handle()

	return cometHandler
}
