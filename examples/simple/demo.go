package main

import (
	"net/http"
	"time"

	"github.com/importcjj/comethandler"
)

// New comet handler.
var comet = comethandler.New()


// Wrapper allowed you do something.
func Wrapper(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	comet.ServeHTTP(rw, r)
}


// Broker to receive a message from manager client
// and try it's best to broadcast it to all clients.
func Broker(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	message := r.FormValue("message")
	comet.Broadcast([]byte(message))
}

func main() {

	http.HandleFunc("/websocket", Wrapper)
	http.HandleFunc("/broker", Broker)

	// Send a tick to clients every minute.
	go func() {
		for now := range time.Tick(1 * time.Second) {
			comet.Broadcast([]byte(now.Format(time.RFC850)))
		}
	}()

	http.ListenAndServe(":8080", nil)
}
