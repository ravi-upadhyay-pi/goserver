package main

import (
	"fasthttptest/config"
	"fasthttptest/handler"
	"fasthttptest/log"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"net"
)

var logger = log.Logger{Id: "main"}
var db *pgx.ConnPool
var redisdb *redis.Client
var mainHandler handler.MainHandler

func main() {
	var err error
	config := config.GetConfig()
	db, err = pgx.NewConnPool(config.SqlConfig)
	if err != nil {
		logger.Fatalf("Error connecting to database: %s", err)
	}
	redisdb = redis.NewClient(&config.RedisOptions)
	defer redisdb.Close()
	mainHandler = handler.NewMainHandler(db, redisdb)
	s := &fasthttp.Server{
		Handler: handle,
		Name:    "go",
	}
	logger.Printf("Starting server on %s", config.Port)
	ln := getListener(config.Port)
	err = s.Serve(ln)
	if err != nil {
		logger.Fatalf("Error starting server: %s", err)
	}
}

func getListener(port string) net.Listener {
	ln, err := net.Listen("tcp4", port)
	if err != nil {
		logger.Fatalf("Could not start listening on port: %s | err: %s", port, err)
	}
	return ln
}

func handle(ctx *fasthttp.RequestCtx) {
	handler.Handle(ctx, mainHandler.Handle)
}
