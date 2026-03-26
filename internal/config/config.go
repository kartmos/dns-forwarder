package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Port           int               `mapstructure:"port"`
	HealthPort     int               `mapstructure:"health_port"`
	Forwarding     map[string]string `mapstructure:"forwarding"`
	TimeoutSeconds int               `mapstructure:"timeout_seconds"`
	Workers        int               `mapstructure:"workers"`
	RateLimitRPS   int               `mapstructure:"rate_limit_rps"`
}

type Store struct {
	mu     sync.RWMutex
	config Config
	viper  *viper.Viper
}

func NewStore(path string) (*Store, error) {
	v := viper.New()
	v.SetConfigFile(path)

	store := &Store{viper: v}
	if err := store.load(); err != nil {
		return nil, err
	}

	store.watch()
	return store, nil
}

func (s *Store) Get() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()

	configCopy := s.config
	configCopy.Forwarding = copyMap(s.config.Forwarding)
	return configCopy
}

func (s *Store) load() error {
	if err := s.viper.ReadInConfig(); err != nil {
		return err
	}

	var cfg Config
	if err := s.viper.Unmarshal(&cfg); err != nil {
		return err
	}

	applyDefaults(&cfg)

	if err := validate(cfg); err != nil {
		return err
	}

	s.mu.Lock()
	s.config = cfg
	s.mu.Unlock()

	return nil
}

func (s *Store) watch() {
	s.viper.OnConfigChange(func(e fsnotify.Event) {
		if e.Op&fsnotify.Write != fsnotify.Write {
			return
		}

		if err := s.load(); err != nil {
			log.Printf("[WARN] failed to reload config: %v", err)
			return
		}

		log.Println("[DONE] config reloaded")
	})

	s.viper.WatchConfig()
}

func validate(cfg Config) error {
	if cfg.Port <= 0 {
		return fmt.Errorf("port must be greater than 0")
	}

	if cfg.HealthPort <= 0 {
		return fmt.Errorf("health_port must be greater than 0")
	}

	if len(cfg.Forwarding) == 0 {
		return fmt.Errorf("forwarding rules are empty")
	}

	return nil
}

func applyDefaults(cfg *Config) {
	if cfg.HealthPort <= 0 {
		cfg.HealthPort = 8080
	}

	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 5
	}

	if cfg.Workers <= 0 {
		cfg.Workers = 50
	}

	if cfg.RateLimitRPS <= 0 {
		cfg.RateLimitRPS = 100
	}
}

func copyMap(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}

	return result
}
