package main

import (
	"time"

	"github.com/cunyat/feeder/internal/feeder"
)

const addr = "127.0.0.1:4000"
const maxConn = 5
const ttl = 60 * time.Second

func main() {
	feeder.New(addr, maxConn, ttl).Start()
}
