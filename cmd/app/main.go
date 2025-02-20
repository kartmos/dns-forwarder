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

	// Запуск UDP-сервера
	addr := net.UDPAddr{Port: Config.Port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("[WARN] UDP-сервер не запущен %s\n", err)
	}
	log.Printf("UDP-сервер запущен на 127.0.0.1:%d", addr.Port)
	for {
		buf := make([]byte, 512)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		log.Printf("[Try] forwarding %v client's request to DNS root Server: %v...\n", clientAddr, Config.Forwarding[clientAddr.IP.String()])
		if err != nil {
			log.Fatalln("[WARN Read UDP]")
		}
		//Обработка каждого запроса в отдельной горутине
		go forwarder.HandleRequest(conn, clientAddr, buf[:n], Config.Forwarding[clientAddr.IP.String()])
	}
}
