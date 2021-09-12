package feeder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
)

type MessageHandler func(context.Context, string)

type Store interface {
	Insert(string)
}

type Application struct {
	ttl   time.Duration
	srv   *server.Server
	store Store

	skus   chan string
	sigint chan os.Signal
}

func New(addr string, maxConn int, store Store, ttl time.Duration) *Application {
	a := &Application{
		ttl:    ttl,
		store:  store,
		skus:   make(chan string, 5),
		sigint: make(chan os.Signal, 1),
	}

	a.srv = server.New(addr, maxConn, a.skus)

	return a
}

func (a *Application) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), a.ttl)
	listenSigInt(a.sigint, cancel)

	// this must not be here, initialize a service that listens chan?
	go func(s Store, skus chan string) {
		for sku := range skus {
			s.Insert(sku)
		}
	}(a.store, a.skus)

	err := a.srv.Start(ctx)
	if err != nil {
		fmt.Printf("Got error from server: %s", err)
	}

	a.Shutdown()
}

func (a *Application) Shutdown() {
	close(a.sigint)
	close(a.skus)
}

func listenSigInt(c chan os.Signal, cancel func()) {
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Got SIGINT, stopping server gracefully...")
		cancel()
	}()
}
