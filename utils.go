package main

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type successResp struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type errorResp struct {
	Message string `json:"message"`
}

// SendErrorResp sends an error response
func SendErrorResp(ctx *fasthttp.RequestCtx, code int, message string) {
	ctx.SetStatusCode(code)
	json.NewEncoder(ctx).Encode(&errorResp{
		Message: message,
	})
}

// SendSuccessResp returns success response
func SendSuccessResp(ctx *fasthttp.RequestCtx, message string, data interface{}) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	json.NewEncoder(ctx).Encode(&successResp{
		Data:    data,
		Message: message,
	})
}
