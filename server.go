package main

import (
	"fasthttptest/config"
	"fasthttptest/handler"
	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	var err error
	var db *pgx.ConnPool
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	logfile, err := os.OpenFile("info.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644); if err != nil {
		log.Fatalf("Could not get logfile: %s", err)
	} else {
		log.SetOutput(io.Writer(logfile))
	}
	config := config.GetConfig()
	db, err = pgx.NewConnPool(config.SqlConfig)
	if err != nil {
		log.Fatalf("herror connecting to database: %s", err)
	}
	mainHandler := handler.NewMainHandler(db)
	s := &fasthttp.Server{
		Handler: mainHandler.Handle,
		Name:    "go",
	}
	log.Printf("Starting server on %s", config.Port)
	ln := getListener(config.Port)
	err = s.Serve(ln)
	if err != nil {
		log.Fatalf("herror starting server: %s", err)
	}
}

func getListener(port string) net.Listener {
	ln, err := net.Listen("tcp4", port)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}
