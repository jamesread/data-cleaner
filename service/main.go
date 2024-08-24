package main

import "github.com/jamesread/data-cleaner/internal/httpservers"
import "github.com/jamesread/data-cleaner/internal/grpc"
import "github.com/jamesread/data-cleaner/internal/config"
import log "github.com/sirupsen/logrus"

func main() {
	log.Infof("data-cleaner")

	config := config.GetConfig()

	log.Infof("config %+v", config)

	go httpservers.Start()

	grpc.Start()
}
