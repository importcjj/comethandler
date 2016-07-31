# comethandler

#### simple demo

*the complete demo is in examples/simple*

```go
package main

import (
	"net/http"
	"time"

	"github.com/importcjj/comethandler"
)

// New comet handler.
var comet = comethandler.New()

// Broker to receive a message from manager client
// and try it's best to broadcast it to all clients.
func Broker(rw http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	comet.Broadcast([]byte(message))
}

func main() {

	http.Handle("/websocket", comet)
	http.HandleFunc("/broker", Broker)

	// Send a tick to clients every minute.
	go func() {
		for now := range time.Tick(1 * time.Second) {
			comet.Broadcast([]byte(now.Format(time.RFC850)))
		}
	}()

	http.ListenAndServe(":8080", nil)
}

```
