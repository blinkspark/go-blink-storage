package main

import (
	"log"

	"github.com/blinkspark/go-blink-storage/server"
)

func main() {
	s, err := server.NewServer("test.key", "ds")
	if err != nil {
		log.Panic(err)
	}
	log.Println(s.Host.Addrs())
	log.Println(s.Host.ID().Pretty())
}
