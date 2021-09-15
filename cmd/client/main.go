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

func main() {
	flag.IntVar(&count, "count", 40, "number of skus to send")
	flag.BoolVar(&terminate, "terminate", false, "send terminate at the end")
	flag.Parse()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))

	for i := 0; i < count; i++ {
		sku := utils.GenerateSKU()
		wg.Add(1)
		sent++
		go sendSKU(sku)
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
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprint(conn, sku)
	if err != nil {
		log.Fatalln("could not write to socket:", err)
	}

	wg.Done()

	_ = conn.Close()
}
