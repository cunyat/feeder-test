package main

import (
	"time"

	"github.com/cunyat/feeder/internal/feeder"
	"github.com/cunyat/feeder/pkg/store"
)

const addr = "127.0.0.1:4000"
const maxConn = 5
const ttl = 60 * time.Second

func main() {
	st := store.New()
	app := feeder.New(addr, maxConn, st, ttl)

	app.Start()
}
