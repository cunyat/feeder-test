package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cunyat/feeder/pkg/utils"
)

var chunks, max int
var terminate bool
var wg sync.WaitGroup

func main() {
	flag.IntVar(&chunks, "chunks", 40, "chunks of skus to send")
	flag.IntVar(&max, "max", 5, "max skus per chunk (random)")
	flag.BoolVar(&terminate, "terminate", false, "send terminate at the end")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))

	for i := 0; i < chunks; i++ {
		skus := utils.GenerateSKUs((rand.Int() % (max)) + 1)
		wg.Add(1)
		go sendSku(skus)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	go func() {
		<-s
		log.Println("Got SIGINT, stopping server...")
		cancel()
	}()

	go func() {
		wg.Wait()
		if terminate {
			wg.Add(1)
			sendSku([]string{"terminate\n"})
		}
		cancel()
	}()

	<-ctx.Done()
}

func sendSku(skus []string) {
	addr, err := net.ResolveTCPAddr("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}
	for _, sku := range skus {
		_, err := fmt.Fprint(conn, sku)
		log.Fatalln("could not write to socket:", err)
	}

	wg.Done()

	// ignore error, message was sent
	_ = conn.Close()
}
