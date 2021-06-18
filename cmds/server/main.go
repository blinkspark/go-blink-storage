package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/blinkspark/go-blink-storage/server"
	"github.com/blinkspark/go-util"
)

func main() {
	StartServer("server.key", "ds", 12233)
}

func StartServer(keyPath, dsPath string, port int) error {
	s, err := server.NewServer(keyPath, dsPath, port)
	if err != nil {
		return err
	}
	annouceServerAddrs(s)
	topic, err := s.PubSub.Join("nealfree.cf/main")
	if err != nil {
		return err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	util.Ignore(sub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			log.Panic(err)
		}
		topic.Publish(context.Background(), raw)
		log.Println("body:", string(raw))
		log.Println(r.PostForm)
	})
	http.ListenAndServe(":10080", nil)

	return nil
}

func annouceServerAddrs(s *server.Server) {
	prettyID := s.Host.ID().Pretty()
	for _, addr := range s.Host.Addrs() {
		fmt.Printf("%s/p2p/%s\n", addr, prettyID)
	}
}
