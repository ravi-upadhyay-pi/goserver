package handler

import (
	"github.com/valyala/fasthttp"
)

const namePrefix = "/name"
const getCountPrefix = "/count"

type plaintext struct {
	count int
}

func (h *plaintext) handle(ctx *fasthttp.RequestCtx) (interface{}, error) {
	h.count += 1
	if isPrefix(ctx, namePrefix){return h.getName(ctx)
	} else if isPrefix(ctx, getCountPrefix){return h.getCount(ctx)
	} else{return notFound(ctx)}
}

func (h *plaintext) getName(ctx *fasthttp.RequestCtx) (string, error) {
	return "test", nil
}

func (h *plaintext) getCount(ctx *fasthttp.RequestCtx) (int, error) {
	return h.count, nil
}