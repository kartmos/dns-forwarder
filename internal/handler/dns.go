package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/kartmos/dns-forwarder.git/internal/metrics"
	"github.com/kartmos/dns-forwarder.git/internal/service"
	"github.com/miekg/dns"
)

type DNSHandler struct {
	forwarder   *service.Forwarder
	rateLimiter *service.RateLimiter
	metrics     *metrics.Metrics
}

func NewDNSHandler(forwarder *service.Forwarder, rateLimiter *service.RateLimiter, metricStore *metrics.Metrics) *DNSHandler {
	return &DNSHandler{
		forwarder:   forwarder,
		rateLimiter: rateLimiter,
		metrics:     metricStore,
	}
}

func (h *DNSHandler) HandleRequest(writer dns.ResponseWriter, request *dns.Msg) {
	start := time.Now()
	network := writer.LocalAddr().Network()

	clientIP, err := clientIPFromAddr(writer.RemoteAddr())
	if err != nil {
		log.Printf("[WARN] failed to parse client address: %v", err)
		h.metrics.IncErrors()
		writeServerFailure(writer, request)
		return
	}

	questionName := "-"
	questionType := uint16(0)
	if len(request.Question) > 0 {
		questionName = request.Question[0].Name
		questionType = request.Question[0].Qtype
	}

	if !h.rateLimiter.Allow(clientIP) {
		h.metrics.IncRateLimited()
		log.Printf("[WARN] rate limit exceeded client=%s network=%s qname=%s qtype=%d", clientIP, network, questionName, questionType)
		writeRefused(writer, request)
		return
	}

	h.metrics.IncRequests(network)
	h.metrics.IncInFlight()
	defer h.metrics.DecInFlight()

	log.Printf("[INFO] request started client=%s network=%s qname=%s qtype=%d", clientIP, network, questionName, questionType)

	response, err := h.forwarder.Resolve(context.Background(), clientIP, request, network)
	if err != nil {
		h.metrics.IncErrors()
		log.Printf("[ERROR] resolve failed client=%s network=%s qname=%s duration=%s err=%v", clientIP, network, questionName, time.Since(start), err)
		writeServerFailure(writer, request)
		return
	}

	if err := writer.WriteMsg(response); err != nil {
		h.metrics.IncErrors()
		log.Printf("[ERROR] failed to write response client=%s network=%s qname=%s err=%v", clientIP, network, questionName, err)
		return
	}

	h.metrics.AddDuration(network, time.Since(start))
	log.Printf("[INFO] request finished client=%s network=%s qname=%s duration=%s", clientIP, network, questionName, time.Since(start))
}

type Server struct {
	udpServer *dns.Server
	tcpServer *dns.Server
}

func NewServer(port int, dnsHandler dns.Handler) *Server {
	address := fmt.Sprintf(":%d", port)

	return &Server{
		udpServer: &dns.Server{
			Addr:    address,
			Net:     "udp",
			Handler: dnsHandler,
		},
		tcpServer: &dns.Server{
			Addr:    address,
			Net:     "tcp",
			Handler: dnsHandler,
		},
	}
}

func (s *Server) Start() error {
	errCh := make(chan error, 2)
	var wg sync.WaitGroup

	startServer := func(server *dns.Server, serverName string) {
		defer wg.Done()
		log.Printf("[DONE] %s server started on %s", serverName, server.Addr)
		if err := server.ListenAndServe(); err != nil {
			errCh <- fmt.Errorf("%s server error: %w", serverName, err)
		}
	}

	wg.Add(2)
	go startServer(s.udpServer, "udp")
	go startServer(s.tcpServer, "tcp")

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.udpServer.ShutdownContext(ctx); err != nil {
		return err
	}

	if err := s.tcpServer.ShutdownContext(ctx); err != nil {
		return err
	}

	return nil
}

func clientIPFromAddr(addr net.Addr) (string, error) {
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}

	return host, nil
}

func writeServerFailure(writer dns.ResponseWriter, request *dns.Msg) {
	response := new(dns.Msg)
	response.SetRcode(request, dns.RcodeServerFailure)

	if err := writer.WriteMsg(response); err != nil {
		log.Printf("[WARN] failed to write servfail: %v", err)
	}
}

func writeRefused(writer dns.ResponseWriter, request *dns.Msg) {
	response := new(dns.Msg)
	response.SetRcode(request, dns.RcodeRefused)

	if err := writer.WriteMsg(response); err != nil {
		log.Printf("[WARN] failed to write refused response: %v", err)
	}
}
