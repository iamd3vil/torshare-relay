package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
	"github.com/valyala/fasthttp"
)

const (
	// DataDir is where Tor instance data is stored
	DataDir = "data"
)

func init() {
	initConfig()
}

func main() {
	hub, err := NewHub()
	if err != nil {
		log.Fatalf("error while initialzing hub: %v", err)
	}

	router := fasthttprouter.New()
	router.POST("/v1/relay", hub.newChannel)
	router.GET("/v1/relay", hub.sendChannelMeta)

	server := &fasthttp.Server{
		Handler: router.Handler,
	}

	// Initialize a tor service
	t, onion, err := newTorService()
	if err != nil {
		log.Fatalf("error while starting tor service: %v", err)
	}

	log.Printf("Starting relay on: %v", onion.ID)

	defer t.Close()
	defer onion.Close()

	log.Fatal(server.Serve(onion))
}

func newTorService() (*tor.Tor, *tor.OnionService, error) {
	var privateKey ed25519.PrivateKey
	_, err := os.Stat(cfg.App.PrivateKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			keyPair, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return nil, nil, err
			}
			privateKey = keyPair.PrivateKey()

			// Write to disk
			f, err := os.Create(cfg.App.PrivateKeyPath)
			if err != nil {
				return nil, nil, err
			}
			_, err = f.Write(privateKey)
			if err != nil {
				return nil, nil, err
			}
			f.Sync()
		} else {
			return nil, nil, err
		}
	}

	privateKey, err = ioutil.ReadFile(cfg.App.PrivateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	t, err := tor.Start(context.Background(), &tor.StartConf{
		DataDir: "/home/sarat/projects/personal/torshare-relay/data/",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while starting Tor service: %v", err)
	}

	onion, err := t.Listen(context.Background(), &tor.ListenConf{
		RemotePorts: []int{80},
		Version3:    true,
		Key:         privateKey,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("error while starting Onion service: %v", err)
	}

	return t, onion, nil
}
