package engine

import (
	"github.com/autom8ter/gonet/grpc/api"
	"github.com/autom8ter/gonet/grpc/config"
	"github.com/autom8ter/gonet/grpc/middleware"
	pbnet "github.com/autom8ter/source/gen/go/util/net"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Option configures a gRPC and a gateway server.
type Option func(*config.GrpcConfig)

func createConfig(opts []Option) *config.GrpcConfig {
	c := config.CreateDefaultConfig()
	for _, f := range opts {
		f(c)
	}
	return c
}

// WithServers returns an Option that sets gRPC service server implementation(s).
func WithServers(svrs ...api.Server) Option {
	return func(c *config.GrpcConfig) {
		c.Servers = append(c.Servers, svrs...)
	}
}

// WithAddr returns an Option that sets an network address for a gRPC and a gateway server.
func WithAddr(network, addr string) Option {
	return func(c *config.GrpcConfig) {
		WithGrpcAddr(network, addr)(c)
		WithGatewayAddr(network, addr)(c)
	}
}

// WithGrpcAddr returns an Option that sets an network address for a gRPC server.
func WithGrpcAddr(network, addr string) Option {
	return func(c *config.GrpcConfig) {
		c.GrpcAddr = &pbnet.Network{
			Network: network,
			Address: addr,
		}
	}
}

// WithGrpcInternalAddr returns an Option that sets an network address connected by a gateway server.
func WithGrpcInternalAddr(network, addr string) Option {
	return func(c *config.GrpcConfig) {
		c.GrpcInternalAddr = &pbnet.Network{
			Network: network,
			Address: addr,
		}
	}
}

// WithGatewayAddr returns an Option that sets an network address for a gateway server.
func WithGatewayAddr(network, addr string) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayAddr = &pbnet.Network{
			Network: network,
			Address: addr,
		}
	}
}

// WithGrpcServerUnaryInterceptors returns an Option that sets unary interceptor(s) for a gRPC server.
func WithGrpcServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(c *config.GrpcConfig) {
		c.GrpcServerUnaryInterceptors = append(c.GrpcServerUnaryInterceptors, interceptors...)
	}
}

// WithGrpcServerStreamInterceptors returns an Option that sets stream interceptor(s) for a gRPC server.
func WithGrpcServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(c *config.GrpcConfig) {
		c.GrpcServerStreamInterceptors = append(c.GrpcServerStreamInterceptors, interceptors...)
	}
}

// WithGatewayServerUnaryInterceptors returns an Option that sets unary interceptor(s) for a gRPC client used by a gateway server.
func WithGatewayServerUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayServerUnaryInterceptors = append(c.GatewayServerUnaryInterceptors, interceptors...)
	}
}

// WithGatewayServerStreamInterceptors returns an Option that sets stream interceptor(s) for a gRPC client used by a gateway server.
func WithGatewayServerStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayServerStreamInterceptors = append(c.GatewayServerStreamInterceptors, interceptors...)
	}
}

// WithGrpcServerOptions returns an Option that sets grpc.ServerOption(s) to a gRPC server.
func WithGrpcServerOptions(opts ...grpc.ServerOption) Option {
	return func(c *config.GrpcConfig) {
		c.GrpcServerOption = append(c.GrpcServerOption, opts...)
	}
}

// WithGatewayDialOptions returns an Option that sets grpc.DialOption(s) to a gRPC clinet used by a gateway server.
func WithGatewayDialOptions(opts ...grpc.DialOption) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayDialOption = append(c.GatewayDialOption, opts...)
	}
}

// WithGatewayMuxOptions returns an Option that sets runtime.ServeMuxOption(s) to a gateway server.
func WithGatewayMuxOptions(opts ...runtime.ServeMuxOption) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayMuxOptions = append(c.GatewayMuxOptions, opts...)
	}
}

// WithGatewayServerMiddlewares returns an Option that sets middleware(s) for http.api.Server to a gateway server.
func WithGatewayServerMiddlewares(middlewares ...middleware.HTTPServerMiddleware) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayServerMiddlewares = append(c.GatewayServerMiddlewares, middlewares...)
	}
}

// WithGatewayServerConfig returns an Option that specifies http.api.Server configuration to a gateway server.
func WithGatewayServerConfig(cfg *config.HTTPServerConfig) Option {
	return func(c *config.GrpcConfig) {
		c.GatewayServerConfig = cfg
	}
}

// WithPassedHeader returns an Option that sets configurations about passed headers for a gateway server.
func WithPassedHeader(decider middleware.PassedHeaderDeciderFunc) Option {
	return WithGatewayServerMiddlewares(middleware.CreatePassingHeaderMiddleware(decider))
}

// WithDefaultLogger returns an Option that sets default grpclogger.LoggerV2 object.
func WithDefaultLogger() Option {
	return func(c *config.GrpcConfig) {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
	}
}