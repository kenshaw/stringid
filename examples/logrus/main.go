// examples/logrus/main.go
package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	goji "goji.io"
	"goji.io/pat"

	"github.com/brankas/stringid"
)

func main() {
	// create logger
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Out = os.Stdout

	// create mux
	mux := goji.NewMux()
	mux.Use(stringid.Middleware())
	mux.HandleFunc(pat.New("/*"), func(res http.ResponseWriter, req *http.Request) {
		logger.WithField("id", stringid.FromRequest(req)).Infof("incoming request for %s", req.URL.Path)
	})

	http.ListenAndServe(":3000", mux)
}
