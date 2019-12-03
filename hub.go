package main

import (
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/sethvargo/go-diceware/diceware"
	"github.com/valyala/fasthttp"
)

const (
	// Expiry is the expiry for keys in Redis
	Expiry = 15 * time.Minute
)

// Hub will hold app environment
type Hub struct {
	rClient *redis.Client
}

// NewHub retuns and new hub instance
func NewHub() (*Hub, error) {
	h := &Hub{
		rClient: redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Addr,
		}),
	}

	return h, nil
}

func (h *Hub) newChannel(ctx *fasthttp.RequestCtx) {
	// Generate a channel name
	list, err := diceware.Generate(6)
	if err != nil {
		log.Printf("couldn't generate channel name: %v", err)
		SendErrorResp(ctx, fasthttp.StatusInternalServerError, "couldn't generate channel name")
		return
	}

	channel := strings.Join(list, "-")

	// Store the body with channel as key in Redis
	_, err = h.rClient.Set(channel, ctx.PostBody(), Expiry).Result()
	if err != nil {
		log.Printf("couldn't store channel name: %v", err)
		SendErrorResp(ctx, fasthttp.StatusInternalServerError, "couldn't store channel name")
		return
	}

	SendSuccessResp(ctx, "", map[string]string{
		"channel": channel,
	})
}

func (h *Hub) sendChannelMeta(ctx *fasthttp.RequestCtx) {
	channel := ctx.QueryArgs().Peek("channel")
	if len(channel) == 0 {
		SendErrorResp(ctx, fasthttp.StatusBadRequest, "channel name can't be empty")
		return
	}

	res, err := h.rClient.Get(string(channel)).Bytes()
	if err != nil {
		SendErrorResp(ctx, fasthttp.StatusBadRequest, "invalid channel name")
		return
	}

	ctx.Write(res)
}
