package forwarder

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

type FrdConfig struct {
	Port       int
	Forwarding map[string]string
}

func HandleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte, dnsServer string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := forwardToDNSServer(ctx, request, dnsServer)
	if err != nil {
		log.Printf("[DO] Cancel request -> %s", err)
		cancel()
		return
	}

	if _, err := conn.WriteToUDP(response, addr); err != nil {
		log.Printf("[DO] Cancel response -> %s", err)
		cancel()
	}
	log.Printf("[Done] Forward Client: %v -> Sever: %v\n", addr, dnsServer)
}

func CheckConfig(config *FrdConfig) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		if e.Op&fsnotify.Write == fsnotify.Write {
			if err := viper.ReadInConfig(); err != nil {
				log.Printf("[WARN] Failed to read config file: %v", err)
				return
			}
			if viper.IsSet("forwarding") {
				config.Forwarding = viper.GetStringMapString("forwarding")
			} else {
				log.Println("[WARN] Key 'forwarding' not found in config")
			}
		}
	})
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("[WARN] Failed to read config file: %v", err)
		return
	}

	config.Port = viper.GetInt("port")
	config.Forwarding = viper.GetStringMapString("forwarding")
}

func forwardToDNSServer(parent context.Context, request []byte, dnsServer string) ([]byte, error) {

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	msg := new(dns.Msg)
	if err := msg.Unpack(request); err != nil {
		log.Printf("[WARN] Unpack request error %s", err)
		return nil, err
	}

	client := new(dns.Client)

	reply, _, err := client.ExchangeContext(ctx, msg, dnsServer)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("[WARN] DNS request timed out")
		} else {
			log.Printf("[ERROR] DNS exchange failed: %v", err)
		}
		return nil, err
	}

	response, err := reply.Pack()
	if err != nil {
		log.Printf("[WARN] Pack response error %s", err)
		return nil, err
	}
	return response, nil
}
