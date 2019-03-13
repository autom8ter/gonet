package gonet

import (
	"fmt"
	"github.com/autom8ter/gonet/config"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	_ "google.golang.org/genproto/googleapis/rpc/errdetails" // Pull in errdetails
	"google.golang.org/grpc"
	"net/http"
	"os"
)

type GrpcGateway struct {
	*Router
	v    *viper.Viper
	port int
}

type RegisterFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type GrpcGatewayConfig struct {
	EnvPrefix    string
	DialOptions  []grpc.DialOption
	RegisterFunc RegisterFunc
}

func NewGrpcGateway(ctx context.Context, cfg *GrpcGatewayConfig, r *Router) *GrpcGateway {
	v := config.SetupViper(cfg.EnvPrefix)
	c := &config.ProxyConfig{
		Endpoint:             v.GetString("endpoint"),
		LogLevel:             v.GetString("log_level"),
		LogHeaders:           v.GetBool("log_headers"),
		CorsAllowOrigin:      v.GetString("cors.allow-origin"),
		CorsAllowCredentials: v.GetString("cors.allow-credentials"),
		CorsAllowMethods:     v.GetString("cors.allow-methods"),
		CorsAllowHeaders:     v.GetString("cors.allow-headers"),
		ApiPrefix:            v.GetString("proxy.api-prefix"),
	}
	gw := config.SetupGateway()
	if err := cfg.RegisterFunc(ctx, gw, c.Endpoint, cfg.DialOptions); err != nil {
		logrus.Fatalf("failed to register grpc gateway from endpoint: %s", err.Error())
	}
	fmt.Printf("registered grpc endpoint:  %s\n", c.Endpoint)
	fmt.Printf("registered gateway handler:  %s\n", c.ApiPrefix)
	r.Mux().Handle(c.ApiPrefix, handlers.CustomLoggingHandler(os.Stdout, http.StripPrefix(c.ApiPrefix[:len(c.ApiPrefix)-1], config.AllowCors(c, gw)), config.LogFormatter(c)))

	return &GrpcGateway{
		Router: r,
		port:   v.GetInt("proxy.port"),
		v:      v,
	}
}
