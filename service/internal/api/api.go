package api

import (
	"net/http"
	"context"
	pb "github.com/jamesread/data-cleaner/gen/data_cleaner/api/v1"
	dcapiv1connect "github.com/jamesread/data-cleaner/gen/data_cleaner/api/v1/dcapiv1connect"
	etlapi "github.com/jamesread/data-cleaner/internal/etlapi"
	"github.com/jamesread/data-cleaner/internal/config"
	"connectrpc.com/connect"
)

type dataCleanerApi struct {
	etl *etlapi.EtlApi
}

func (s *dataCleanerApi) Import(ctx context.Context, in *connect.Request[pb.ImportRequest]) (*connect.Response[pb.ImportResponse], error) {
	res := s.etl.Import()

	return connect.NewResponse(res), nil
}

func (s *dataCleanerApi) Export(ctx context.Context, in *connect.Request[pb.ExportRequest]) (*connect.Response[pb.ExportResponse], error) {
	if in.Msg.RunImport {
		s.etl.Import()
	}

	ret := &pb.ExportResponse{
	}

	return connect.NewResponse(ret), nil
}

func (s *dataCleanerApi) Reload(ctx context.Context, in *connect.Request[pb.ReloadRequest]) (*connect.Response[pb.ReloadResponse], error) {
	config.ReloadConfig()

	res := &pb.ReloadResponse{}

	return connect.NewResponse(res), nil
}

func (s *dataCleanerApi) Load(ctx context.Context, in *connect.Request[pb.LoadRequest]) (*connect.Response[pb.LoadResponse], error) {
	s.etl.Load()

	res := &pb.LoadResponse{
	}

	return connect.NewResponse(res), nil
}

func NewServer() *dataCleanerApi {
	server := &dataCleanerApi{}
	server.etl = etlapi.NewEtlApi()
	return server
}

func GetNewHandler() (string, http.Handler) {
	server := NewServer()

	path, handler := dcapiv1connect.NewDataCleanerServiceHandler(server)

	return path, handler
}
