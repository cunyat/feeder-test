package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cunyat/feeder/pkg/utils"
)

var count, sent int
var terminate bool
var wg sync.WaitGroup

const workers = 8

func main() {
	flag.IntVar(&count, "count", 40, "number of skus to send")
	flag.BoolVar(&terminate, "terminate", false, "send terminate at the end")
	flag.Parse()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))

	skus := make(chan string, 8)

	for i := 0; i < workers; i++ {
		go func(skus chan string) {
			for sku := range skus {
				sendSKU(sku)
			}
		}(skus)
	}

	for i := 0; i < count; i++ {
		sku := utils.GenerateSKU()
		wg.Add(1)
		sent++
		skus <- sku
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
			sendSKU("terminate\n")
		}
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("total skus:", sent)
}

func sendSKU(sku string) {
	addr, err := net.ResolveTCPAddr("tcp", ":4000")
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fmt.Fprint(conn, sku)
	if err != nil {
		log.Fatalln("could not write to socket:", err)
	}

	wg.Done()

	_ = conn.Close()
}
