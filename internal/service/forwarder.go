package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kartmos/dns-forwarder.git/internal/config"
	"github.com/miekg/dns"
)

type Forwarder struct {
	configStore *config.Store
	sem         chan struct{}
}

func NewForwarder(configStore *config.Store) *Forwarder {
	cfg := configStore.Get()

	return &Forwarder{
		configStore: configStore,
		sem:         make(chan struct{}, cfg.Workers),
	}
}

func (f *Forwarder) Resolve(ctx context.Context, clientIP string, request *dns.Msg, network string) (*dns.Msg, error) {
	cfg := f.configStore.Get()
	upstream, ok := cfg.Forwarding[clientIP]
	if !ok {
		return nil, fmt.Errorf("upstream not found for client %s", clientIP)
	}

	f.sem <- struct{}{}
	defer func() {
		<-f.sem
	}()

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second

	client := &dns.Client{
		Net:     network,
		Timeout: timeout,
	}

	reply, _, err := client.ExchangeContext(ctx, request, upstream)
	if err != nil {
		log.Printf("[ERROR] upstream request failed client=%s network=%s upstream=%s err=%v", clientIP, network, upstream, err)
		return nil, err
	}

	if network == "udp" && reply.Truncated {
		log.Printf("[INFO] truncated udp response client=%s upstream=%s, retry with tcp", clientIP, upstream)
		tcpClient := &dns.Client{
			Net:     "tcp",
			Timeout: timeout,
		}

		reply, _, err = tcpClient.ExchangeContext(ctx, request, upstream)
		if err != nil {
			log.Printf("[ERROR] tcp fallback failed client=%s upstream=%s err=%v", clientIP, upstream, err)
			return nil, err
		}
	}

	log.Printf("[INFO] upstream request finished client=%s network=%s upstream=%s", clientIP, network, upstream)
	return reply, nil
}
