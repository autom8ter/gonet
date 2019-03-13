package engine

import (
	"github.com/autom8ter/gonet/grpc/api"
	"github.com/autom8ter/gonet/grpc/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"net"
)

// GrpcServer wraps grpc.Server setup process.
type GrpcServer struct {
	server *grpc.Server
	*config.GrpcConfig
}

// NewGrpcServer creates GrpcServer instance.
func NewGrpcServer(c *config.GrpcConfig) api.Interface {
	s := grpc.NewServer(c.ServerOptions()...)
	reflection.Register(s)
	for _, svr := range c.Servers {
		svr.RegisterWithServer(s)
	}
	return &GrpcServer{
		server:     s,
		GrpcConfig: c,
	}
}

// Serve implements Server.Shutdown
func (s *GrpcServer) Serve(l net.Listener) error {
	grpclog.Infof("gRPC server is starting %v", l.Addr())

	err := s.server.Serve(l)

	grpclog.Infof("gRPC server stopped: %v", err)

	return errors.Wrap(err, "failed to serve gRPC server")
}

// Shutdown implements Server.Shutdown
func (s *GrpcServer) Shutdown() {
	s.server.GracefulStop()
}
