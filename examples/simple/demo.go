package main

import (
	"net/http"
	"time"

	"github.com/importcjj/comethandler"
)

// New comet handler.
var comet = comethandler.New()

func Manager(rw http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	comet.Broadcast([]byte(message))
}

func main() {

	http.Handle("/websocket", comet)
	http.HandleFunc("/manager", Manager)

	go func() {
		for now := range time.Tick(1 * time.Second) {
			comet.Broadcast([]byte(now.Format(time.RFC850)))
		}
	}()

	http.ListenAndServe(":8080", nil)
}
