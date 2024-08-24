package httpservers

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/jamesread/data-cleaner/gen/grpc"
	"github.com/jamesread/data-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var cfg *config.Config

func Start() {
	cfg = config.GetConfig()

	go startRestApiServer()

	go startSingleFrontend()
}

func startRestApiServer() {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	pb.RegisterDataCleanerServiceHandlerFromEndpoint(context.Background(), mux, cfg.Network.BindGrpc, opts)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startSingleFrontend() {
	apiUrl, _ := url.Parse("http://" + cfg.Network.BindRest)
	apiProxy := httputil.NewSingleHostReverseProxy(apiUrl)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		apiProxy.ServeHTTP(w, r)
	})
	mux.Handle("/", http.FileServer(http.Dir("../frontend/")))
	srv := &http.Server{
		Addr:    cfg.Network.BindProxy,
		Handler: mux,
	}

	log.Infof("Starting proxy on %s", cfg.Network.BindProxy)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
