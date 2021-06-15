package main

import (
	"fmt"
	"log"

	"github.com/blinkspark/go-blink-storage/server"
)

func main() {
	s, err := server.NewServer("test.key", "ds")
	if err != nil {
		log.Panic(err)
	}
	id := s.Host.ID().Pretty()
	for _, addr := range s.Host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr, id)
	}

	defer s.Host.Close()
}
