package grpc

import (
	"context"
	pb "github.com/jamesread/data-cleaner/gen/grpc"
	"github.com/jamesread/data-cleaner/internal/api"
	"github.com/jamesread/data-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"net"
)

type dataCleanerApi struct {
}

func (s *dataCleanerApi) Import(ctx context.Context, in *pb.ImportRequest) (*pb.ImportResponse, error) {
	res := api.Import()

	return res, nil
}

func (s *dataCleanerApi) Export(ctx context.Context, in *pb.ExportRequest) (*httpbody.HttpBody, error) {
	csvdata := api.Export()

	ret := &httpbody.HttpBody{
		ContentType: "text/plain",
		Data:        csvdata,
	}

	return ret, nil
}

func (s *dataCleanerApi) Reload(ctx context.Context, in *pb.ReloadRequest) (*pb.ReloadResponse, error) {
	config.ReloadConfig()

	return &pb.ReloadResponse{}, nil
}

func newServer() *dataCleanerApi {
	server := &dataCleanerApi{}
	return server
}

func Start() {
	log.Infof("Starting grpc server on " + config.GetConfig().Network.BindGrpc)

	lis, err := net.Listen("tcp", config.GetConfig().Network.BindGrpc)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDataCleanerServiceServer(grpcServer, newServer())
	err = grpcServer.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
