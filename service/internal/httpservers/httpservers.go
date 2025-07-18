package httpservers

import (
	"github.com/jamesread/data-cleaner/internal/frontend"
	"github.com/jamesread/data-cleaner/internal/api"
	log "github.com/sirupsen/logrus"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Start() {
	mux := http.NewServeMux()

	apipath, apihandler := api.GetNewHandler()

	log.Infof("API path: /api/%s", apipath)

	mux.Handle("/api"+apipath, http.StripPrefix("/api", apihandler))
	mux.Handle("/", http.StripPrefix("/", frontend.GetNewHandler()))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
