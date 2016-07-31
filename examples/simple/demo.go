package main

import (
	"net/http"

	"fmt"
	"io/ioutil"
	"os"

	"github.com/importcjj/comethandler"
)

func ReadHTML(name string) []byte {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}

	HTML, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return HTML
}

var ClientHTML, ManagerHTML []byte

func init() {
	ClientHTML = ReadHTML("client.html")
	ManagerHTML = ReadHTML("notification.html")
}

// New comet handler.
var comet = comethandler.New()

func Client(rw http.ResponseWriter, r *http.Request) {
	rw.Write(ClientHTML)
}

func Manager(rw http.ResponseWriter, r *http.Request) {
	rw.Write(ManagerHTML)
}

func SendFunc(rw http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	fmt.Println(message)
	comet.Broadcast([]byte(message))
}

func main() {
	http.HandleFunc("/client", Client)
	http.HandleFunc("/manager", Manager)

	http.Handle("/websocket", comet)
	http.HandleFunc("/send", SendFunc)

	// go func() {
	// 	for now := range time.Tick(5 * time.Second) {
	// 		comet.Broadcast([]byte(now.Format(time.RFC850)))
	// 	}
	// }()

	http.ListenAndServe(":8080", nil)
}
