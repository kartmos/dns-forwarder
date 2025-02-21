package main

import (
	"log"
	"net"

	"github.com/kartmos/dns-forwarder.git/internal/forwarder"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config/config.yaml")
}

func main() {

	var Config forwarder.FrdConfig
	forwarder.CheckConfig(&Config)
	log.Println("[DONE] Set options from config file")

	addr := net.UDPAddr{Port: Config.Port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("[WARN] Failed to start server %v: %s\n", addr.AddrPort(), err)
	}
	log.Printf("UDP-server start on 127.0.0.1:%d", addr.Port)
	for {
		buf := make([]byte, 512)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln("[WARN] Err with read UDP request from client")
		}
		log.Printf("[Try] forwarding %v client's request to DNS root Server: %v...\n", clientAddr, Config.Forwarding[clientAddr.IP.String()])
		go forwarder.HandleRequest(conn, clientAddr, buf[:n], Config.Forwarding[clientAddr.IP.String()])
	}
}
