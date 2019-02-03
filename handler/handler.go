package handler

import (
	"encoding/json"
	"fasthttptest/ferror"
	"fasthttptest/log"
	"github.com/valyala/fasthttp"
	"strings"
)

func notFound(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return nil, ferror.New(404, ferror.NotFound, "404 Not Found")
}

type handler func (logger log.Logger, ctx *fasthttp.RequestCtx) (interface{}, error)

func handle(logger log.Logger, ctx *fasthttp.RequestCtx, handler handler) {
	res, err := handler(logger, ctx)
	if err != nil {
		var herr *ferror.Error
		ctx.Response.Reset()
		switch v := err.(type) {
		case *ferror.Error:
			herr = v
		default:
			herr = ferror.New(500, ferror.InternalServerError, v.Error())
		}
		json, _ := json.Marshal(herr)
		ctx.SetStatusCode(herr.HttpStatusCode)
		ctx.SetBody(json)
	} else {
		json, _ := json.Marshal(res)
		ctx.SetBody(json)
	}
}

func Handle(ctx *fasthttp.RequestCtx, handler handler) {
	path := strings.Split(string(ctx.Path()), "/")[1:]
	ctx.SetUserValue("relativePath", path)
	logger := getLogger(ctx)
	handle(logger, ctx, handler)
}

func isPrefix(ctx *fasthttp.RequestCtx, prefix ...string) bool {
	path := ctx.UserValue("relativePath").([]string)
	if len(prefix) > len(path) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if prefix[i][0] == ':' {
			ctx.SetUserValue(prefix[i], path[i])
		} else if path[i] != prefix[i] {
			return false
		}
	}
	ctx.SetUserValue("relativePath", path[len(prefix):])
	return true
}

func match(ctx *fasthttp.RequestCtx, method string, prefix ...string) bool {
	return string(ctx.Method()) == method && isPrefix(ctx, prefix...) &&
		len(ctx.UserValue("relativePath").([]string)) == 0
}

func getLogger(ctx *fasthttp.RequestCtx) log.Logger {
	var requestId string
	requestIdBytes := ctx.Request.Header.Peek("X-Request-Id")
	if len(requestIdBytes) == 0 {
		requestId = "unknown"
	} else {
		requestId = string(requestIdBytes)
	}
	return log.Logger{Id: requestId}
}
