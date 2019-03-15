package config

import (
	"encoding/json"
	"github.com/autom8ter/goproxyrpc/pkg/errors"
	"github.com/autom8ter/goproxyrpc/pkg/health"
	"github.com/autom8ter/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type ProxyConfig struct {
	// The Backend gRPC service to listen to.
	Endpoint string
	// The log level to use
	LogLevel string
	// Whether to log request headers
	LogHeaders bool
	// Value to set for Access-Control-Allow-Origin header.
	CorsAllowOrigin string
	// Value to set for Access-Control-Allow-Credentials header.
	CorsAllowCredentials string
	// Value to set for Access-Control-Allow-Methods header.
	CorsAllowMethods string
	// Value to set for Access-Control-Allow-Headers header.
	CorsAllowHeaders string
	// Prefix that this gateway is running on. For example, if your API endpoint
	// was "/foo/bar" in your protofile, and you wanted to run APIs under "/api",
	// set this to "/api/".
	ApiPrefix string
}

func AllowCors(cfg *ProxyConfig, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		CorsAllowOrigin := cfg.CorsAllowOrigin
		if CorsAllowOrigin == "*" {
			if origin := req.Header.Get("Origin"); origin != "" {
				CorsAllowOrigin = origin
			}
		}
		w.Header().Set("Access-Control-Allow-Origin", CorsAllowOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", cfg.CorsAllowCredentials)
		w.Header().Set("Access-Control-Allow-Methods", cfg.CorsAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", cfg.CorsAllowHeaders)
		if req.Method == "OPTIONS" && req.Header.Get("Access-Control-Request-Method") != "" {
			return
		}
		handler.ServeHTTP(w, req)
	})
}

func LogFormatter(cfg *ProxyConfig) handlers.LogFormatter {

	// Setup logrus
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
	})
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

	return func(writer io.Writer, params handlers.LogFormatterParams) {

		host, _, err := net.SplitHostPort(params.Request.RemoteAddr)
		if err != nil {
			host = params.Request.RemoteAddr
		}

		uri := params.Request.RequestURI

		// Requests using the CONNECT method over HTTP/2.0 must use
		// the authority field (aka r.Host) to identify the target.
		// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
		if params.Request.ProtoMajor == 2 && params.Request.Method == "CONNECT" {
			uri = params.Request.Host
		}
		if uri == "" {
			uri = params.URL.RequestURI()
		}

		duration := int64(time.Now().Sub(params.TimeStamp) / time.Millisecond)

		fields := logrus.Fields{
			"host":       host,
			"url":        uri,
			"duration":   duration,
			"status":     params.StatusCode,
			"method":     params.Request.Method,
			"request":    params.Request.RequestURI,
			"remote":     params.Request.RemoteAddr,
			"size":       params.Size,
			"referer":    params.Request.Referer(),
			"user_agent": params.Request.UserAgent(),
			"request_id": params.Request.Header.Get("x-request-id"),
		}

		// Only append headers if explicitly enabled
		if cfg.LogHeaders {
			if headers, err := json.Marshal(params.Request.Header); err == nil {
				fields["headers"] = string(headers)
			} else {
				fields["header_error"] = err.Error()
			}
		}

		logrus.WithFields(fields).WithTime(params.TimeStamp).Infof("%s %s %d", params.Request.Method, uri, params.StatusCode)
	}
}

// SetupViper returns a viper configuration object
func SetupViper(envPrefix string) *viper.Viper {
	util.DotEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if envPrefix != "" {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("proxy.port", 8080)
	viper.SetDefault("proxy.api-prefix", "/")
	viper.SetDefault("endpoint", "http://localhost:3000")

	err := viper.ReadInConfig()
	if err != nil {
		errors.New("Could not read config", err).FailIfErr()
	}
	endPoint := viper.GetString("endpoint")
	if viper.InConfig("proxy.api-prefix") {
		viper.Set("proxy.api-prefix", util.SanitizeApiPrefix(viper.GetString("proxy.api-prefix")))
	}
	if endPoint == "" {
		errors.New("", errors.NewErr("please provide a non-empty endpoint in your configuration")).FailIfErr()
	}
	errors.New("failed to ping grpc endpoint", health.New(endPoint).Once().Do()).FailIfErr()
	return viper.GetViper()
}

type Auth0Config struct {
	Domain       string   `json:"domain"`
	Audience     string   `json:"audience"`
	Callback     string   `json:"callback"`
	Scopes       []string `json:"scopes"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	cfg          *oauth2.Config
	Cookies      *sessions.CookieStore
}
