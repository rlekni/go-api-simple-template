package server

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	meter               = otel.Meter("tea-blends-api")
	httpRequestsTotal   metric.Int64Counter
	httpRequestDuration metric.Float64Histogram
)

func init() {
	var err error
	httpRequestsTotal, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests."),
	)
	if err != nil {
		slog.Error("failed to create http_requests_total counter", "error", err)
		os.Exit(1)
	}

	httpRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests."),
		metric.WithUnit("s"),
	)
	if err != nil {
		slog.Error("failed to create http_request_duration_seconds histogram", "error", err)
		os.Exit(1)
	}
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip tracking for metrics, statsviz and pprof endpoints
		if r.URL.Path == "/metrics" ||
			(len(r.URL.Path) >= 11 && r.URL.Path[:11] == "/metrics-ui") ||
			(len(r.URL.Path) >= 12 && r.URL.Path[:12] == "/debug/pprof") {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()

		// Ensure path doesn't have high cardinality by using the route pattern
		path := chi.RouteContext(r.Context()).RoutePattern()
		if path == "" {
			path = "unmapped"
		}

		status := strconv.Itoa(ww.Status())

		attrs := metric.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", path),
			attribute.String("status", status),
		)

		httpRequestsTotal.Add(r.Context(), 1, attrs)
		httpRequestDuration.Record(r.Context(), duration, attrs)
	})
}

// FreeMemoryHandler manually triggers GC and releases memory to the OS.
func FreeMemoryHandler(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	w.Write([]byte("Memory released to OS"))
}
