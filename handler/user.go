package handler

import (
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"time"
)

const testPrefix = "/test"

type user struct {
	db *pgx.ConnPool
}

func (uh *user) handle(ctx *fasthttp.RequestCtx) (interface{}, error) {
	if isPrefix(ctx, testPrefix){return uh.test(ctx)
	} else{return notFound(ctx)}
}

func (uh *user) test(ctx *fasthttp.RequestCtx) (interface{}, error) {
	var result time.Time
	err := uh.db.QueryRow("select now()").Scan(&result)
	return result, err
}