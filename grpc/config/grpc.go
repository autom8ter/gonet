package config

import (
	"github.com/autom8ter/gonet/grpc/api"
	"github.com/autom8ter/gonet/grpc/middleware"
	netpb "github.com/autom8ter/source/gen/go/util/net"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

// Config contains configurations of gRPC and Gateway server.
type GrpcConfig struct {
	GrpcAddr                        *netpb.Network
	GrpcInternalAddr                *netpb.Network
	GatewayAddr                     *netpb.Network
	Servers                         []api.Server
	GrpcServerUnaryInterceptors     []grpc.UnaryServerInterceptor
	GrpcServerStreamInterceptors    []grpc.StreamServerInterceptor
	GatewayServerUnaryInterceptors  []grpc.UnaryClientInterceptor
	GatewayServerStreamInterceptors []grpc.StreamClientInterceptor
	GrpcServerOption                []grpc.ServerOption
	GatewayDialOption               []grpc.DialOption
	GatewayMuxOptions               []runtime.ServeMuxOption
	GatewayServerConfig             *HTTPServerConfig
	MaxConcurrentStreams            uint32
	GatewayServerMiddlewares        []middleware.HTTPServerMiddleware
}

func CreateDefaultConfig() *GrpcConfig {
	grpcaddr := os.Getenv("GRPC_PORT")
	gatewayaddr := os.Getenv("GATEWAY_PORT")
	if grpcaddr == "" {
		grpcaddr = "0.0.0.0:3000"
	}
	if gatewayaddr == "" {
		gatewayaddr = "0.0.0.0:8080"
	}
	config := &GrpcConfig{
		GrpcAddr:         nil,
		GrpcInternalAddr: netpb.NewNetwork(grpcaddr, "tcp"),
		GatewayAddr:      netpb.NewNetwork(gatewayaddr, "tcp"),
		GatewayServerConfig: &HTTPServerConfig{
			ReadTimeout:  8 * time.Second,
			WriteTimeout: 8 * time.Second,
			IdleTimeout:  2 * time.Minute,
		},
		MaxConcurrentStreams:     1000,
		GatewayServerMiddlewares: nil,
	}
	return config
}

func (c *GrpcConfig) ServerOptions() []grpc.ServerOption {
	return append(
		[]grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(c.GrpcServerUnaryInterceptors...),
			grpc_middleware.WithStreamServerChain(c.GrpcServerStreamInterceptors...),
			grpc.MaxConcurrentStreams(c.MaxConcurrentStreams),
		},
		c.GrpcServerOption...,
	)
}

func (c *GrpcConfig) ClientOptions() []grpc.DialOption {
	return append(
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
				return net.Dial(c.GrpcInternalAddr.Network, a)
			}),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(c.GatewayServerUnaryInterceptors...),
			),
			grpc.WithStreamInterceptor(
				grpc_middleware.ChainStreamClient(c.GatewayServerStreamInterceptors...),
			),
		},
		c.GatewayDialOption...,
	)
}
