package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blinkspark/go-blink-storage/server"
)

func main() {
	s, err := server.NewServer("server.key", "ds")
	if err != nil {
		log.Panic(err)
	}
	id := s.Host.ID().Pretty()
	for _, addr := range s.Host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr, id)
	}

	topic, err := s.PubSub.Join("test")
	if err != nil {
		log.Panic(err)
	}

	sub, err := topic.Subscribe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		ctx := context.Background()
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				log.Println(err)
				continue
			}
			raw := msg.Data
			from := msg.GetFrom()
			log.Printf("from:%s,msg:%s\n", from.Pretty(), string(raw))
			log.Println(topic.ListPeers())
			log.Println(s.Host.Peerstore().Peers())

		}
	}()
	select {}
	// defer s.Host.Close()
}
