package service

import (
	"sync"
	"time"
)

type clientRate struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	limit   int
	window  time.Duration
	clients map[string]*clientRate
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]*clientRate),
	}
}

func (r *RateLimiter) Allow(clientIP string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	client, ok := r.clients[clientIP]
	if !ok {
		r.clients[clientIP] = &clientRate{
			count:       1,
			windowStart: now,
		}
		return true
	}

	if now.Sub(client.windowStart) >= r.window {
		client.count = 1
		client.windowStart = now
		return true
	}

	if client.count >= r.limit {
		return false
	}

	client.count++
	return true
}
