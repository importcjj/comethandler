package main

import (
	"net/http"
	"runtime"

	"time"

	"github.com/importcjj/comethandler"
)

// New comet handler.
var comet = comethandler.New()

func MyHandlerFunc(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	comet.ServeHTTP(rw, r)
}

func ManagerFunc(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	message := r.FormValue("message")
	comet.Broadcast([]byte(message))
}

func main() {
	runtime.GOMAXPROCS(2)
	http.HandleFunc("/websocket", MyHandlerFunc)
	http.HandleFunc("/manager", ManagerFunc)

	go func() {
		for now := range time.Tick(3 * time.Second) {
			comet.Broadcast([]byte(now.Format(time.RFC850)))
		}
	}()

	http.ListenAndServe(":8080", nil)
}
