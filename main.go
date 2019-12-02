package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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

	log.Fatal(fasthttp.ListenAndServe(cfg.App.Addr, router.Handler))
}
