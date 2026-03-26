package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kartmos/dns-forwarder.git/internal/config"
	"github.com/kartmos/dns-forwarder.git/internal/metrics"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(port int, configStore *config.Store, metricStore *metrics.Metrics) *HTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler())
	mux.HandleFunc("/readyz", readyHandler(configStore))
	mux.HandleFunc("/metrics", metricsHandler(metricStore))

	return &HTTPServer{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *HTTPServer) Start() error {
	log.Printf("[DONE] http server started on %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func healthHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("ok"))
	}
}

func readyHandler(configStore *config.Store) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		cfg := configStore.Get()
		if len(cfg.Forwarding) == 0 {
			writer.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(writer).Encode(map[string]string{"status": "not ready"})
			return
		}

		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(map[string]string{"status": "ready"})
	}
}

func metricsHandler(metricStore *metrics.Metrics) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; version=0.0.4")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte(metricStore.ExportPrometheus()))
	}
}
