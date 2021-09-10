package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/store"
)

const addr = "127.0.0.1:4000"
const maxConn = 5

var ttl = 60 * time.Second

func main() {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not listen in port %d: %s", 4000, err.Error())
	}

	skus := make(chan string, 1)
	store := store.New()
	go func() {
		for sku := range skus {
			store.Insert(sku)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), ttl)
	for i := 0; i < maxConn; i++ {
		go listen(ctx, ln, skus)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		fmt.Println("Got SIGINT, stopping server gracefully...")
		cancel()
	}()

	<-ctx.Done()
	close(c)
	close(skus)

	fmt.Printf("Received %d unique product skus, %d duplicates, %d discard values\n", store.SKUCount(), store.DuplicatedCount(), 0)
}

func listen(ctx context.Context, ln net.Listener, out chan string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("stopping listener...")
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				log.Fatalf("error accepting a new connection: %s", err.Error())
			}

			for {
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println("could not read incomming message: ", err.Error())
				}

				out <- msg
			}
		}
	}
}
