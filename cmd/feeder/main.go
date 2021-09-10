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
	"strings"
	"time"
)

const addr = "127.0.0.1:4000"
const maxConn = 5

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not listen in port %d: %s", 4000, err.Error())
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))

	go func() {
		<-c
		log.Println("Got SIGINT, stopping server...")
		cancel()
	}()

	for i := 0; i < maxConn; i++ {
		go listen(ctx, ln)
	}

	<-ctx.Done()

}

func listen(ctx context.Context, ln net.Listener) {
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

			var skus []string
			for {
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Println("could not read incomming message: ", err.Error())
				}

				msg = strings.TrimRight(msg, "\n")
				skus = append(skus, msg)
			}

			fmt.Printf("got: %s\n", skus)
		}
	}
}
