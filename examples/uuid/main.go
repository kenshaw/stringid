// examples/uuid/main.go
package main

import (
	"fmt"
	"net/http"

	goji "goji.io"
	"goji.io/pat"

	"github.com/brankas/stringid"
)

func main() {
	mux := goji.NewMux()
	mux.Use(stringid.Middleware(
		stringid.WithGenerator(stringid.NewUUIDGenerator()),
	))
	mux.HandleFunc(pat.New("/*"), func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "request id: %s\n", stringid.FromContext(req.Context()))
	})
	http.ListenAndServe(":3000", mux)
}
