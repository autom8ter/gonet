package gonet

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/autom8ter/gonet/db"
	"github.com/autom8ter/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"net/http"
	"net/http/pprof"
	"os"
)

func init() {
	fs = afero.NewOsFs()
}

var fs afero.Fs

type Router struct {
	fs     *afero.HttpFs
	addr   string
	router *mux.Router
	chain  *negroni.Negroni
	db     *db.MongoDB
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(rw, req)
}

func NewRouter(addr string) *Router {
	m := mux.NewRouter()
	n := negroni.Classic()
	httpFs := afero.NewHttpFs(fs)
	return &Router{
		fs:     httpFs,
		addr:   addr,
		router: m,
		db:     nil,
		chain:  n,
	}
}

func NewMongoRouter(addr, colName, connectionStr, databaseName string) *Router {
	return &Router{
		addr:   addr,
		router: mux.NewRouter(),
		chain:  negroni.Classic(),
		db:     db.NewMongoDB(colName, connectionStr, databaseName),
	}
}

func (r *Router) Mongo() *db.MongoDB {
	if r.db == nil {
		panic("Database uninitialized, us NewMongoRouter to add a database connection")
	}
	return r.db
}

func (r *Router) SwitchMongo(m *db.MongoDB) {
	r.db = m
}

func (r *Router) SwitchNegroni(n *negroni.Negroni) {
	r.chain = n
}
func (r *Router) SwitchAddr(a string) {
	r.addr = a
}
func (r *Router) SwitchRouter(router *mux.Router) {
	r.router = router
}

func (r *Router) WithDebug() {
	WithDebug(r.router)
}

func (r *Router) WithStaticViews() {
	WithStaticViews(r.router)
}

func (r *Router) WithMetrics() {
	WithMetrics(r.router)
}

func (r *Router) WithJWT(signingKey string, debug bool, path string, handler http.Handler) {
	WithJWT(signingKey, debug, path, handler, r.router)
}

func (r *Router) Serve() {
	r.chain.UseHandler(r.router)
	r.chain.Run(r.addr)
}

func (r *Router) NotImplememntedFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Write([]byte("Not Implemented"))
	}
}

func (r *Router) OnErrorUnauthorized(w http.ResponseWriter, req *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

func (r *Router) OnErrorInternal(w http.ResponseWriter, req *http.Request, err string) {
	http.Error(w, err, http.StatusInternalServerError)
}

func (r *Router) GenerateJWT(signingKey string, claims map[string]interface{}) (string, error) {
	return util.GenerateJWT(signingKey, claims)
}

func (r *Router) SetResponseHeaders(headers map[string]string, w http.ResponseWriter) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}

func (r *Router) GetHeader(key string, w http.ResponseWriter) string {
	return w.Header().Get(key)
}

func (r *Router) DelHeader(key string, w http.ResponseWriter) {
	w.Header().Del(key)
}

func (r *Router) Do(r2 *http.Request) (*http.Response, error) {
	client := http.DefaultClient
	return client.Do(r2)
}

func (r *Router) DoClient(r2 *http.Request, client *http.Client) (*http.Response, error) {
	return client.Do(r2)
}

func (r *Router) Stringify(obj interface{}) string {
	return util.ToPrettyJsonString(obj)
}

func (r *Router) JSONify(obj interface{}) []byte {
	return util.ToPrettyJson(obj)
}

func (r *Router) RandomTokenString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (r *Router) RandomToken(length int) []byte {
	b := make([]byte, length)
	rand.Read(b)
	return b
}

func (r *Router) DerivePassword(counter uint32, password_type, password, user, site string) string {
	return util.DerivePassword(counter, password, password, user, site)
}

func (r *Router) GeneratePrivateKey(typ string) string {
	return util.GeneratePrivateKey(typ)
}

func (r *Router) Render(s string, data interface{}) string {
	return util.Render(s, data)
}

func (r *Router) SetSessionValFunc(cookieStore *sessions.CookieStore, name string, vals map[string]interface{}) http.HandlerFunc {
	return SetSessionValFunc(cookieStore, name, vals)
}

func (r *Router) NewSessionFSStore(path, key string) *sessions.FilesystemStore {
	return NewSessionFileStore(path, key)
}

func (r *Router) NewSessionCookieStore(key string) *sessions.CookieStore {
	return NewSessionCookieStore(key)
}

func (r *Router) AddFlashSessionFunc(cookieStore *sessions.CookieStore, name string, val interface{}, vars ...string) http.HandlerFunc {
	return AddFlashSessionFunc(cookieStore, name, val, vars...)
}

func (r *Router) Router() *mux.Router {
	if r.router == nil {
		return mux.NewRouter()
	}
	return r.router
}

func (r *Router) HTTPFS() *afero.HttpFs {
	if r.fs == nil {
		return afero.NewHttpFs(fs)
	}
	return r.fs
}

func (r *Router) RequestVars(req *http.Request) map[string]string {
	return RequestVars(req)
}

func JWTMiddleware(singingKey string, debug bool) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return singingKey, nil
		},
		Debug:         debug,
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func WithJWT(signingKey string, debug bool, path string, handler http.Handler, r *mux.Router) {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		},
		Debug: debug,
		Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader,
			jwtmiddleware.FromParameter("code"),
			jwtmiddleware.FromParameter("auth-code")),
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.Handle(path, negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(handler),
	))
}

func WithDebug(r *mux.Router) {
	r.HandleFunc("/debug", func(w http.ResponseWriter, request *http.Request) {
		fmt.Println("registered handler: ", "/debug")
		w.Write([]byte(fmt.Sprintln("Status: 💡 API is up and running 💡 ")))
		w.Write([]byte(fmt.Sprintln("---------------------------------------------------------------------")))
		w.Write([]byte(fmt.Sprintln("Configuration Settings:")))
		w.Write([]byte(fmt.Sprintln(util.ToPrettyJsonString(viper.AllSettings()))))
		w.Write([]byte(fmt.Sprintln("---------------------------------------------------------------------")))
		w.Write([]byte(fmt.Sprintln("Environmental Variables:")))
		w.Write([]byte(fmt.Sprintln(os.Environ())))
		w.Write([]byte(fmt.Sprintln("---------------------------------------------------------------------")))
		w.Write([]byte(fmt.Sprintln("Registered Router Paths:")))
		err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			type routeLog struct {
				Name         string
				PathRegExp   string
				PathTemplate string
				HostTemplate string
				Methods      []string
			}
			meth, _ := route.GetMethods()
			host, _ := route.GetHostTemplate()
			pathreg, _ := route.GetPathRegexp()
			pathtemp, _ := route.GetPathTemplate()
			rout := &routeLog{
				Name:         route.GetName(),
				PathRegExp:   pathreg,
				PathTemplate: pathtemp,
				HostTemplate: host,
				Methods:      meth,
			}
			w.Write([]byte(util.ToPrettyJson(rout)))
			return nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(fmt.Sprintln("---------------------------------------------------------------------")))
	})
	fmt.Println("registered handler: ", "/debug/pprof/")
	r.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
	fmt.Println("registered handler: ", "/debug/pprof/cmdline")
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	fmt.Println("registered handler: ", "/debug/pprof/profile")
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	fmt.Println("registered handler: ", "/debug/pprof/symbol")
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	fmt.Println("registered handler: ", "/debug/pprof/trace")
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
}

func WithStaticViews(r *mux.Router) {
	// On the default page we will simply serve our static index page.
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	fmt.Println("registered file server handler: ", "./views/")
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	fmt.Println("registered file server handler: ", "./static/")
}

func WithMetrics(r *mux.Router) {
	var (
		inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "in_flight_requests",
			Help: "A gauge of requests currently being served by the wrapped handler.",
		})

		counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "api_requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"code", "method"},
		)

		// duration is partitioned by the HTTP method and handler. It uses custom
		// buckets based on the expected request duration.
		duration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests.",
				Buckets: []float64{.025, .05, .1, .25, .5, 1},
			},
			[]string{"handler", "method"},
		)

		// responseSize has no labels, making it a zero-dimensional
		// ObserverVec.
		responseSize = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "response_size_bytes",
				Help:    "A histogram of response sizes for requests.",
				Buckets: []float64{200, 500, 900, 1500},
			},
			[]string{},
		)
	)

	// Register all of the metrics in the standard registry.
	prometheus.MustRegister(inFlightGauge, counter, duration, responseSize)
	var chain http.Handler
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pth, _ := route.GetPathTemplate()
		chain = promhttp.InstrumentHandlerInFlight(inFlightGauge,
			promhttp.InstrumentHandlerDuration(duration.MustCurryWith(prometheus.Labels{"handler": pth}),
				promhttp.InstrumentHandlerCounter(counter,
					promhttp.InstrumentHandlerResponseSize(responseSize, route.GetHandler()),
				),
			),
		)
		route = route.Handler(chain)
		return nil
	})
	fmt.Println("registered handler: ", "/metrics")
	r.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
}

func RequestVars(req *http.Request) map[string]string {
	return mux.Vars(req)
}
