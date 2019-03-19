package gonet

import (
	"context"
	"github.com/autom8ter/gonet/pkg/config"
	"github.com/autom8ter/gonet/pkg/netutil"
	"github.com/autom8ter/gonet/pkg/router"
	"github.com/autom8ter/util"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"google.golang.org/grpc"
	"net/http"
	"net/http/httputil"
	"strings"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h(rw, r)
}

func (h HandlerFunc) AsHandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		h(writer, request)
	}
}

func (h HandlerFunc) Before(after http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		after(w, r)
	}
}

func (h HandlerFunc) After(before http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		before(w, r)
		h(w, r)
	}
}

func (h HandlerFunc) Chain(chained ...http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		for _, fn := range chained {
			fn(w, r)
		}
	}
}

func (h HandlerFunc) SwitchGRPC(grpcServer *grpc.Server) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func (h HandlerFunc) SwitchJSON(handler http.HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "json") {
			handler.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}



func AsHandlerFunc(handlerFunc http.HandlerFunc) HandlerFunc {
	return HandlerFunc(handlerFunc)
}

type GRPCRouter func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

func (g GRPCRouter) Handler(ctx context.Context, cfg *config.ProxyConfig, opts []grpc.DialOption, fns ...router.RouterFunc) http.Handler {
	m := mux.NewRouter()
	gw := netutil.SetupGateway()
	for _, f := range fns {
		f(m)
	}
	if err := g(ctx, gw, cfg.Endpoint, opts); err != nil {
		logrus.Fatalf("failed to register grpc gateway from endpoint: %s", err.Error())
	}
	m.Handle(cfg.ApiPrefix, gw)
	n := negroni.Classic()
	n.UseHandler(m)
	return n
}

func (g GRPCRouter) HandlerFromConfig(ctx context.Context, envPrefix string, opts []grpc.DialOption, fns ...router.RouterFunc) http.Handler {
	v := config.SetupViper(envPrefix)
	m := mux.NewRouter()
	gw := netutil.SetupGateway()
	for _, f := range fns {
		f(m)
	}
	cfg := &config.ProxyConfig{
		Endpoint:             v.GetString("endpoint"),
		LogLevel:             v.GetString("log_level"),
		LogHeaders:           v.GetBool("log_headers"),
		CorsAllowOrigin:      v.GetString("cors.allow-origin"),
		CorsAllowCredentials: v.GetString("cors.allow-credentials"),
		CorsAllowMethods:     v.GetString("cors.allow-methods"),
		CorsAllowHeaders:     v.GetString("cors.allow-headers"),
		ApiPrefix:            v.GetString("proxy.api-prefix"),
	}
	if err := g(ctx, gw, cfg.Endpoint, opts); err != nil {
		logrus.Fatalf("failed to register grpc gateway from endpoint: %s", err.Error())
	}
	m.Handle(cfg.ApiPrefix, gw)
	n := negroni.Classic()
	n.UseHandler(m)
	return n
}

//ReverseProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler
type ReverseRouter map[string]*httputil.ReverseProxy

//ResponseMiddleware is a function used to modify the response of a reverse proxy
type ResponseMiddleware func(func(response *http.Response) error) func(response *http.Response) error

//RequestMiddleware is a function used to modify the incoming request of a reverse proxy from a client
type RequestMiddleware func(func(req *http.Request)) func(req *http.Request)

//TransportMiddleware is a function used to modify the http RoundTripper that is used by a reverse proxy. The default RoundTripper is initially http.DefaultTransport
type TransportMiddleware func(tripper http.RoundTripper) http.RoundTripper

//ProxyConfig is used to configure GoProxies reverse proxies
type ProxyConfig struct {
	PathPrefix string
	TargetUrl  string
	Username   string
	Password   string
	Method     string
	Headers    map[string]string
	Form       map[string]string
}

//NewReverseProxy registers a new reverseproxy for each provided ProxyConfig
func NewReverseHandler(configs ...*ProxyConfig) ReverseRouter {
	g := make(map[string]*httputil.ReverseProxy)
	for _, c := range configs {
		g[c.PathPrefix] = &httputil.ReverseProxy{
			Director: util.ProxyRequestFunc(c.TargetUrl, c.Method, c.Username, c.Password, c.Headers, c.Form),
		}
	}
	return g
}

//ModifyResponses takes a Response Middleware function, traverses each registered reverse proxy, and modifies the http response it sends to the client
func (g ReverseRouter) ModifyResponses(middleware ResponseMiddleware) {
	for _, prox := range g {
		prox.ModifyResponse = middleware(prox.ModifyResponse)
	}
}

//ModifyResponses takes a Request Middleware function, traverses each registered reverse proxy, and modifies the http request it sends to its target prior to sending
func (g ReverseRouter) ModifyRequests(middleware RequestMiddleware) {
	for _, prox := range g {
		prox.Director = middleware(prox.Director)
	}
}

//ModifyResponses takes a Transport Middleware function, traverses each registered reverse proxy, and modifies the http roundtripper it uses
func (g ReverseRouter) ModifyTransport(middleware TransportMiddleware) {
	for _, prox := range g {
		prox.Transport = middleware(prox.Transport)
	}
}

//GetProxy returns the reverse proxy with the registered prefix
func (g ReverseRouter) Get(prefix string) *httputil.ReverseProxy {
	return g[prefix]
}

//AsHandlerFunc converts a ReverseProxy to an http.HandlerFunc for convenience
func (g ReverseRouter) Handler(fns ...router.RouterFunc) http.Handler {
	m := mux.NewRouter()
	for _, f := range fns {
		f(m)
	}
	for path, f := range g {
		m.Handle(path, f)
	}
	n := negroni.Classic()
	n.UseHandler(m)
	return n
}

type Router map[string]HandlerFunc

func (r Router) Handler(fns ...router.RouterFunc) http.Handler {
	m := mux.NewRouter()
	for _, f := range fns {
		f(m)
	}
	for path, f := range r {
		logrus.Debugln("registered handler: ", path)
		m.HandleFunc(path, f)
	}
	n := negroni.Classic()
	n.UseHandler(m)
	return n
}

type Runner struct {
	Addr    string
	Routers []func(router Router)
	Muxers  []router.RouterFunc
}

func NewRunner(addr string, routers []func(router Router), muxers []router.RouterFunc) *Runner {
	return &Runner{Addr: addr, Routers: routers, Muxers: muxers}
}

func (r *Runner) Run() error {
	var rmux = Router{}
	for _, f := range r.Routers {
		f(rmux)
	}
	return http.ListenAndServe(r.Addr, rmux.Handler(r.Muxers...))
}

func GrpcAsHandlerFunc(s *grpc.Server) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.ServeHTTP(w, r)
	}
}