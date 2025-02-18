package main

import (
	"log"
	"net"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config/config.yaml")
}

func main() {

	var config Forward
	CheckConfig(&config)

	// Запуск UDP-сервера
	addr := net.UDPAddr{Port: config.Port}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("[WARN] UDP-сервер не запущен %s\n", err)
	}
	log.Printf("UDP-сервер запущен на 127.0.0.1:%d", addr.Port)
	for {
		buf := make([]byte, 512)
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatalln("[WARN Read UDP]")
		}
		//Обработка каждого запроса в отдельной горутине
		go HandleRequest(conn, clientAddr, buf[:n], config.Forwarding[clientAddr.IP.String()])
	}
}
