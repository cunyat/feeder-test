package main

import (
	"fmt"
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

	fmt.Printf(
		"Received %d unique product skus, %d duplicates, %d discard values\n",
		st.UniqueCount(),
		st.DuplicatedCount(),
		0,
	)
}
