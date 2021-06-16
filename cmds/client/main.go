package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/blinkspark/go-blink-storage/client"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	keyName := flag.String("k", "client.key", "-k client.key")
	flag.Parse()

	c, err := client.NewClient(*keyName)
	if err != nil {
		log.Panic(err)
	}
	log.Println(c.ID().Pretty())
	log.Println(c.Host.Addrs())
	ma, err := multiaddr.NewMultiaddr("/ip4/192.168.0.110/udp/12233/quic/p2p/12D3KooWDbkYQnpfEj3dnkBxdJJFZo2h2xDQG3xJaMaRik8LuboC")
	if err != nil {
		log.Panic(err)
	}
	pi, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		log.Panic(err)
	}
	c.Connect(context.Background(), *pi)

	topic, err := c.PubSub.Join("test")
	if err != nil {
		log.Panic(err)
	}
	sub, err := topic.Subscribe()
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
		}
	}()
	go func() {
		ctx := context.Background()
		for {
			topic.Publish(ctx, []byte("Hello"))
			log.Println(topic.ListPeers())
			log.Println(c.Host.Peerstore().Peers())
			id, err := peer.IDFromString("12D3KooWBGhYRzKrCFjWy21ioD12JR7X2SismkcSyPLZyyBMVPMM")
			if err != nil {
				log.Panic(err)
			}
			pi, err := c.IpfsDHT.FindPeer(ctx, id)
			if err != nil {
				log.Panic(err)
			}
			log.Println("find", pi)
			time.Sleep(time.Second)
		}
	}()
	select {}
}
