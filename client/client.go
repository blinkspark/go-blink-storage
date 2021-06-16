package client

import (
	"context"
	"os"
	"time"

	myutil "github.com/blinkspark/go-util"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
)

type Client struct {
	crypto.PrivKey
	host.Host
	*dht.IpfsDHT
	*pubsub.PubSub
}

func NewClient(keyPath string) (s *Client, err error) {
	ctx := context.Background()
	s = &Client{}
	// gen priv
	if myutil.PathExists(keyPath) {
		privRaw, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		s.PrivKey, err = crypto.UnmarshalPrivateKey(privRaw)
		if err != nil {
			return nil, err
		}
	} else {
		s.PrivKey, _, err = crypto.GenerateKeyPair(crypto.Ed25519, -1)
		if err != nil {
			return nil, err
		}
		privRaw, err := crypto.MarshalPrivateKey(s.PrivKey)
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(keyPath, privRaw, 0666)
		if err != nil {
			return nil, err
		}
	}
	id := libp2p.Identity(s.PrivKey)

	h, err := libp2p.New(ctx,
		id, libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",      // regular tcp connections
			"/ip6/::/tcp/0",           // regular tcp6 connections
			"/ip4/0.0.0.0/udp/0/quic", // a UDP endpoint for the QUIC transport
			"/ip6/::/udp/0/quic",      // a UDP6 endpoint for the QUIC transport
		),
		libp2p.Transport(libp2pquic.NewTransport),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr.NewConnManager(
			50,          // Lowwater
			300,         // HighWater,
			time.Minute, // GracePeriod
		)),
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			s.IpfsDHT, err = dht.New(ctx, h)
			return s.IpfsDHT, err
		}),
		libp2p.EnableAutoRelay(),
	)

	s.PubSub, err = pubsub.NewGossipSub(ctx, h, pubsub.WithPeerExchange(true))
	if err != nil {
		return nil, err
	}

	s.Host = h
	return s, nil
}
