
```go
package main

import (
	"net/http"

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
    // http.HandleFunc("/websocket", comet.Func)
	http.HandleFunc("/manager", Manager)
	http.ListenAndServe(":8080", nil)
}
```
