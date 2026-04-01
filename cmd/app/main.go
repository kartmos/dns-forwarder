package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/kartmos/dns-forwarder.git/internal/config"
	"github.com/kartmos/dns-forwarder.git/internal/handler"
	"github.com/kartmos/dns-forwarder.git/internal/metrics"
	"github.com/kartmos/dns-forwarder.git/internal/service"
	"github.com/miekg/dns"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configStore, err := config.NewStore("config/config.yaml")
	if err != nil {
		log.Fatalf("[WARN] failed to load config: %v", err)
	}
	log.Println("[DONE] config loaded")

	forwarderService := service.NewForwarder(configStore)
	cfg := configStore.Get()
	rateLimiter := service.NewRateLimiter(cfg.RateLimitRPS, time.Second)
	metricStore := metrics.New()
	dnsHandler := handler.NewDNSHandler(forwarderService, rateLimiter, metricStore)

	server := handler.NewServer(cfg.Port, dns.HandlerFunc(dnsHandler.HandleRequest))
	httpServer := handler.NewHTTPServer(cfg.HealthPort, configStore, metricStore)

	go func() {
		if err := httpServer.Start(); err != nil {
			log.Printf("[WARN] http server stopped: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Println("[DONE] shutdown signal received")
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("[WARN] shutdown error: %v", err)
		}
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("[WARN] http shutdown error: %v", err)
		}
	}()

	if err := server.Start(); err != nil {
		log.Fatalf("[WARN] server stopped: %v", err)
	}
}
