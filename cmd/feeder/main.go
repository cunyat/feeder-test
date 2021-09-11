package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
	"github.com/cunyat/feeder/pkg/store"
)

const addr = "127.0.0.1:4000"
const maxConn = 5

var ttl = 60 * time.Second

func main() {
	skus := make(chan string, 1)
	store := store.New()
	go func() {
		for sku := range skus {
			store.Insert(sku)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), ttl)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Got SIGINT, stopping server gracefully...")
		cancel()
	}()

	srv := server.New(addr, maxConn, skus)
	err := srv.Start(ctx)
	if err != nil {
		fmt.Printf("Got error from server: %s", err)
	}

	close(c)
	close(skus)

	fmt.Printf(
		"Received %d unique product skus, %d duplicates, %d discard values\n",
		store.SKUCount(),
		store.DuplicatedCount(),
		0,
	)
}
