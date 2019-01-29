package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
)

const plaintextPrefix = "/plaintext"
const userPrefix = "/user"

type MainHandler struct {
	uh user
	ph plaintext
}

func NewMainHandler(db *pgx.ConnPool) MainHandler {
	return MainHandler {
		uh: user{db},
		ph: plaintext{},
	}
}

func (mh *MainHandler) Handle(ctx *fasthttp.RequestCtx) {
	putCLogger(ctx)
	path := string(ctx.Path())
	ctx.SetUserValue("relativePath", path)
	ctx.SetContentType("application/json")
	if isPrefix(ctx, plaintextPrefix){handle(ctx, mh.ph.handle)
	}else if isPrefix(ctx, userPrefix){handle(ctx, mh.uh.handle)
	}else{handle(ctx, notFound)}
}

func handle(ctx *fasthttp.RequestCtx, handler handler) {
	res, err := handler(ctx)
	if err != nil {
		var herr herror
		ctx.Response.Reset()
		switch v := err.(type) {
		case herror: herr = v
		default: herr = herror{500, "internal server error", v.Error()}
		}
		json, _ := json.Marshal(herr)
		ctx.SetStatusCode(herr.HttpStatusCode)
		ctx.SetBody(json)
	} else {
		json, _ := json.Marshal(res)
		ctx.SetBody(json)
	}
}

func isPrefix(ctx *fasthttp.RequestCtx, prefix string) bool {
	path := relativePath(ctx)
	if strings.HasPrefix(path, prefix) {
		ctx.SetUserValue("relativePath", path[len(prefix):])
		return true
	}
	return false
}

func relativePath(ctx *fasthttp.RequestCtx) string {
	return ctx.UserValue("relativePath").(string)
}

func notFound(ctx *fasthttp.RequestCtx) (interface{}, error){
	return nil, &herror{404, "404 Not Found", "at " + relativePath(ctx)}
}

type herror struct {
	HttpStatusCode int
	ErrorCode string
	Message string
}
func (herr herror) Error() string {return herr.Message}
type handler func (ctx *fasthttp.RequestCtx) (interface{}, error)


type cLog struct {
	cid string
}

func putCLogger(ctx *fasthttp.RequestCtx) cLog {
	var cid string
	cidBytes := ctx.Request.Header.Peek("X-Request-Id")
	if len(cidBytes) == 0 {
		cid = "unknown"
	} else {
		cid = string(cidBytes)
	}
	clogger := cLog{cid: cid}
	ctx.SetUserValue("clogger", clogger)
	return clogger

}

func cLogger(ctx *fasthttp.RequestCtx) cLog {
	return ctx.UserValue("clogger").(cLog)
}

func (clog *cLog) Printf(format string, v ...interface{}) {
	format = clog.cid + ": " + format
	log.Output(2,fmt.Sprintf(format, v...))
}
