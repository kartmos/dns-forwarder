package forwarder

import (
	"context"
	"log"
	"net"
	"time"
	"github.com/kartmos/dns-forwarder.git/cmd/app"
	"github.com/fsnotify/fsnotify"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

func HandleRequest(conn *net.UDPConn, addr *net.UDPAddr, request []byte, dnsServer string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Перенаправление запроса на DNS-сервер
	response, err := ForwardToDNSServer(ctx, request, dnsServer)
	if err != nil {
		log.Printf("[DO] Cancel request -> %s", err)
		cancel()
		return
	}

	// Отправка ответа клиенту
	conn.WriteToUDP(response, addr)
	log.Printf("[Done] Forward Client: %v -> Sever: %v\n\n", addr, dnsServer)
	return
}

func CheckConfig(config *app.Forward) {
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

	// Чтение конфигурации с помощью Viper
	config.Port = viper.GetInt("port")
	config.Forwarding = viper.GetStringMapString("forwarding")
}

func ForwardToDNSServer(parent context.Context, request []byte, dnsServer string) ([]byte, error) {

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
