package grpc

import (
	"context"
	"github.com/autom8ter/gonet/grpc/api"
	"github.com/autom8ter/gonet/grpc/engine"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

func ServeGateway(gwaddr, grpcaddr string, gateway *Gateway) error {
	return NewServerEngine(gwaddr, grpcaddr, gateway).Serve()
}

func ServeAPIGateway(gwaddr, grpcaddr string, gateway *APIGateway) error {
	return NewGatewayEngine(gateway).Serve()
}

func NewServerEngine(gwAddr, grpcAddr string, servers ...api.Server) *engine.Engine {
	return engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			servers...,
		),
		engine.WithGatewayAddr("tcp", gwAddr),
		engine.WithAddr("tcp", grpcAddr),
	)
}

func NewGatewayEngine(gateway api.Gateway) *engine.Engine {
	return engine.New(
		engine.WithDefaultLogger(),
		engine.WithServers(
			gateway.GetServers()...,
		),
		engine.WithGatewayAddr(gateway.GetGatewayNetAddr()),
		engine.WithAddr(gateway.GetGRPCNetAddr()),
	)
}

type Gateway struct {
	RegisterServerFunc  func(*grpc.Server)
	RegisterHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

func NeGateway(serverfn func(*grpc.Server), gwfn func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error) *Gateway {
	return &Gateway{
		RegisterServerFunc:  serverfn,
		RegisterHandlerFunc: gwfn,
	}
}

func (d *Gateway) RegisterWithServer(s *grpc.Server) {
	d.RegisterServerFunc(s)
}

func (d *Gateway) RegisterWithHandler(ctx context.Context, m *runtime.ServeMux, cc *grpc.ClientConn) error {
	return d.RegisterHandlerFunc(ctx, m, cc)
}

func (d *Gateway) AsServer() api.Server {
	return d
}

type APIGateway struct {
	GatewayNetwork string
	GatewayAddr    string
	GRPCNetwork    string
	GRPCAddr       string
	GrpcGateways   []*Gateway
}

func NewAPIGateway(gatewayNetwork string, gatewayAddr string, GRPCNetwork string, GRPCAddr string, grpcGateways ...*Gateway) *APIGateway {
	if len(grpcGateways) == 0 {
		panic("must provide at least one grpc gateway")
	}
	return &APIGateway{GatewayNetwork: gatewayNetwork, GatewayAddr: gatewayAddr, GRPCNetwork: GRPCNetwork, GRPCAddr: GRPCAddr, GrpcGateways: grpcGateways}
}

func (a *APIGateway) AsGateway() api.Gateway {
	return a
}

func (a *APIGateway) GetServers() []api.Server {
	servers := []api.Server{}
	for _, s := range a.GrpcGateways {
		servers = append(servers, s)
	}
	return servers
}

func (a *APIGateway) GetGRPCNetAddr() (string, string) {
	return a.GRPCNetwork, a.GRPCAddr
}

func (a *APIGateway) GetGatewayNetAddr() (string, string) {
	return a.GatewayNetwork, a.GatewayAddr
}
