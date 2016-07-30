package comethandler

import (
	"net/http"
)

// Context is the HTTP request keeper.
type Context struct {
	ReturnSignal   chan bool
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// Wait to receive a message.
// then it can be return
func (c *Context) Wait() {
	<-c.ReturnSignal
}

// Write to the HTTP request's ResponseWriter
// and then let's the request return.
func (c *Context) Write(body []byte) {
	c.ResponseWriter.Write(body)
	c.ReturnSignal <- true
}

// NewContext return a new context.
func NewContext(responseWriter http.ResponseWriter, request *http.Request) *Context {
	return &Context{
		ReturnSignal:   make(chan bool),
		ResponseWriter: responseWriter,
		Request:        request,
	}
}
