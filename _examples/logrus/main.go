// examples/logrus/main.go
package main

import (
	"net/http"
	"os"

	"github.com/brankas/goji"
	"github.com/brankas/stringid"
	"github.com/sirupsen/logrus"
)

func main() {
	// create logger
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Out = os.Stdout

	// create mux
	mux := goji.New()
	mux.Use(stringid.Middleware())
	mux.HandleFunc(goji.NewPathSpec("/*"), func(res http.ResponseWriter, req *http.Request) {
		logger.WithField("id", stringid.FromRequest(req)).Infof("incoming request for %s", req.URL.Path)
	})

	http.ListenAndServe(":3000", mux)
}
