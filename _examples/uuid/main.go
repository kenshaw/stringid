// examples/uuid/main.go
package main

import (
	"fmt"
	"net/http"

	"github.com/brankas/goji"
	"github.com/brankas/stringid"
)

func main() {
	mux := goji.New()
	mux.Use(stringid.Middleware(
		stringid.WithGenerator(stringid.NewUUIDGenerator()),
	))
	mux.HandleFunc(goji.NewPathSpec("/*"), func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "request id: %s\n", stringid.FromContext(req.Context()))
	})
	http.ListenAndServe(":3000", mux)
}
