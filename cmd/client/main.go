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
		skus := generateSKUs(int(randInt(max-1)) + 1)
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
		sendSku([]string{"terminate"})

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
		fmt.Fprintf(conn, sku)
	}

	conn.Close()
	wg.Done()
}

var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var numbers = []byte("0123456789")

func generateSKUs(count int) []string {
	skus := make([]string, count)
	for i := range skus {
		skus[i] = fmt.Sprintf("%s-%s\n", pick(letters, 4), pick(numbers, 4))
	}

	return skus
}

func pick(source []byte, count int) string {
	b := make([]byte, count)

	for i := range b {
		b[i] = source[randInt(len(source))]
	}

	return string(b)
}

func randInt(max int) int64 {
	return rand.NewSource(time.Now().UnixNano()).Int63() % int64(max)
}
