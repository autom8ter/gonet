package api

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type Gateway interface {
	AsGateway() Gateway
	GetServers() []Server
	GetGRPCNetAddr() (string, string)
	GetGatewayNetAddr() (string, string)
}

// Server is an interface for representing gRPC server implementations.
type Server interface {
	AsServer() Server
	RegisterWithServer(*grpc.Server)
	RegisterWithHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

type Interface interface {
	Serve(l net.Listener) error
	Shutdown()
}
