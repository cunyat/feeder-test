package main

import (
	"flag"
	"time"

	"github.com/cunyat/feeder/internal/feeder"
)

const addr = "127.0.0.1:4000"
const maxConn = 5
const ttl = 60 * time.Second

func main() {
	outfile := flag.String("outfile", "out/skus.txt", "output file to write skus")
	flag.Parse()

	app := feeder.New(addr, maxConn, ttl, *outfile)
	app.Start()
}
