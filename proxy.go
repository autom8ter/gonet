package gonet

import (
	"github.com/autom8ter/util"
	"net/http"
	"net/http/httputil"
)

//ReverseProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler
type ReverseProxy struct {
	*Router
	proxies map[string]*httputil.ReverseProxy
}

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
func NewReverseProxy(addr string, configs ...*ProxyConfig) *ReverseProxy {
	g := &ReverseProxy{
		Router:  NewRouter(addr),
		proxies: make(map[string]*httputil.ReverseProxy),
	}
	for _, c := range configs {
		g.proxies[c.PathPrefix] = &httputil.ReverseProxy{
			Director: util.ProxyRequestFunc(c.TargetUrl, c.Method, c.Username, c.Password, c.Headers, c.Form),
		}
	}
	for path, prox := range g.proxies {
		g.Router.Router().Handle(path, prox)
	}
	return g
}

//ModifyResponses takes a Response Middleware function, traverses each registered reverse proxy, and modifies the http response it sends to the client
func (g *ReverseProxy) ModifyResponses(middleware ResponseMiddleware) {
	for _, prox := range g.proxies {
		prox.ModifyResponse = middleware(prox.ModifyResponse)
	}
}

//ModifyResponses takes a Request Middleware function, traverses each registered reverse proxy, and modifies the http request it sends to its target prior to sending
func (g *ReverseProxy) ModifyRequests(middleware RequestMiddleware) {
	for _, prox := range g.proxies {
		prox.Director = middleware(prox.Director)
	}
}

//ModifyResponses takes a Transport Middleware function, traverses each registered reverse proxy, and modifies the http roundtripper it uses
func (g *ReverseProxy) ModifyTransport(middleware TransportMiddleware) {
	for _, prox := range g.proxies {
		prox.Transport = middleware(prox.Transport)
	}
}

//Proxies returns all registered reverse proxies as a map of prefix:reverse proxy
func (g *ReverseProxy) Proxies() map[string]*httputil.ReverseProxy {
	return g.proxies
}

//GetProxy returns the reverse proxy with the registered prefix
func (g *ReverseProxy) GetProxy(prefix string) *httputil.ReverseProxy {
	return g.proxies[prefix]
}

//AsHandlerFunc converts a ReverseProxy to an http.HandlerFunc for convenience
func (g *ReverseProxy) AsHandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		g.ServeHTTP(writer, request)
	}
}
