package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/arl/statsviz"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"go-api-simple-template/internal/database"
	"go-api-simple-template/internal/handler"
	"go-api-simple-template/internal/service"
)

type Server struct {
	port int

	db *pgxpool.Pool

	teaBlendHandler *handler.TeaBlendHandler
	teaBlendService *service.TeaBlendService

	Router chi.Router
	API    huma.API
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080 // Default port
	}

	dbPool, _ := database.InitPgxPool()

	// Services
	tbs := service.NewTeaBlendService(dbPool)

	// Handlers
	tbh := handler.NewTeaBlendHandler(tbs)

	router := chi.NewRouter()

	// Add a custom logger to see the protocol clearly in stdout
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.DebugContext(r.Context(), "received request",
				slog.String("proto", r.Proto),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path))
			next.ServeHTTP(w, r)
		})
	})

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(MetricsMiddleware)
	router.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "tea-blends-api")
	})

	// Profiling and Monitoring

	// Statsviz UI
	srv, _ := statsviz.NewServer(statsviz.Root("/metrics-ui"))
	router.Get("/metrics-ui/ws", srv.Ws())
	router.Handle("/metrics-ui/*", srv.Index())
	router.Get("/metrics-ui", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics-ui/", http.StatusMovedPermanently)
	})

	router.Route("/debug/pprof", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
		r.Handle("/allocs", pprof.Handler("allocs"))
		r.Handle("/block", pprof.Handler("block"))
		r.Handle("/goroutine", pprof.Handler("goroutine"))
		r.Handle("/heap", pprof.Handler("heap"))
		r.Handle("/mutex", pprof.Handler("mutex"))
		r.Handle("/threadcreate", pprof.Handler("threadcreate"))
	})

	config := huma.DefaultConfig("Tea Blends API", "1.0.0")
	config.CreateHooks = nil
	config.DocsRenderer = huma.DocsRendererScalar
	config.DocsPath = "/docs"
	config.Servers = []*huma.Server{
		{URL: fmt.Sprintf("http://localhost:%d", port), Description: "Local Server"},
	}

	api := humachi.New(router, config)

	s := &Server{
		port: port,

		db:              dbPool,
		teaBlendService: tbs,
		teaBlendHandler: tbh,

		Router: router,
		API:    api,
	}

	s.RegisterRoutes()

	h2s := &http2.Server{
		MaxConcurrentStreams: 250,
		IdleTimeout:          10 * time.Second,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      h2c.NewHandler(router, h2s),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
