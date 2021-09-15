package feeder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
)

type MessageHandler interface {
	HandleMessage(string)
}

type Store interface {
	Insert(string)
}

type Application struct {
	ttl     time.Duration
	srv     *server.Server
	handler MessageHandler

	skus   chan string
	sigint chan os.Signal
}

func New(addr string, maxConn int, store Store, ttl time.Duration) *Application {
	a := &Application{
		ttl:    ttl,
		handler: NewManager(store, ValidateSKU),
		skus:   make(chan string, 5),
		sigint: make(chan os.Signal, 1),
	}

	a.srv = server.New(addr, maxConn, a.skus)

	return a
}

func (a *Application) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), a.ttl)

	listenSigInt(a.sigint, cancel)

	go func(handler MessageHandler, skus chan string, cancel context.CancelFunc) {
		for sku := range skus {
			if IsTerminateSequence(sku) {
				cancel()
			}

			handler.HandleMessage(sku)
		}
	}(a.handler, a.skus, cancel)

	err := a.srv.Start(ctx)
	if err != nil {
		fmt.Printf("Got error from server: %s", err)
	}

	a.Shutdown()

	fmt.Println()
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
