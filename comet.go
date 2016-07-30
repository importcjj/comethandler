package comethandler

import "net/http"

type CometHandler struct {
	MessageBox chan []byte
	Pool       *ContextPool
}

func (c *CometHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	context := NewContext(rw, r)
	err := c.Pool.Put(context)
	if err != nil {
		rw.Write([]byte("Connection reset by peer"))
		return
	}
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
		clientNum := c.Pool.Length
		for clientNum > 0 {
			requestContext, err := c.Pool.Get()
			if err != nil {
				break
			}
			clientNum--
			requestContext.Write(message)
		}
	}
}

func New() *CometHandler {
	cometHandler := new(CometHandler)
	cometHandler.MessageBox = make(chan []byte, 10)
	cometHandler.Pool = NewContextPool()

	go cometHandler.handle()

	return cometHandler
}
