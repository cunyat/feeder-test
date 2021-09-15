package feeder

import (
	"context"
	"fmt"
	"github.com/cunyat/feeder/pkg/store"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
)

type Application struct {
	// time to live of the server
	ttl     time.Duration
	// server that listens for incoming messages
	srv     *server.Server
	// store to keep skus deduplicated
	store   *store.DeduplicatedStore
	// manager that handles incoming messages
	handler *Manager
	// writter writes results in a file
	writter *OutputWritter

	skus   chan string
	sigint chan os.Signal
}

func New(addr string, maxConn int, store *store.DeduplicatedStore, ttl time.Duration) *Application {
	a := &Application{
		ttl:    ttl,
		store: store,
		handler: NewManager(store, ValidateSKU),
		writter: NewOutputWritter(fmt.Sprintf("out/skus-%d.txt", time.Now().Unix())),
		skus:   make(chan string, 5),
		sigint: make(chan os.Signal, 1),
	}

	a.srv = server.New(addr, maxConn, a.skus)

	return a
}

func (a *Application) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), a.ttl)

	listenSigInt(a.sigint, cancel)

	go func(handler *Manager, skus chan string, cancel context.CancelFunc) {
		for sku := range skus {
			if IsTerminateSequence(sku) {
				fmt.Println("Got terminate sequece, stopping server gracefully...")
				cancel()
				return
			}

			a.handler.HandleMessage(sku)
		}
	}(a.handler, a.skus, cancel)

	err := a.srv.Start(ctx)
	if err != nil {
		fmt.Printf("Got error from server: %s", err)
	}

	a.Shutdown()

	if err := a.writter.Write(a.store.GetReader()); err != nil {
		fmt.Println("Error writting output file: ", err)
	}

	fmt.Printf(
		"Received %d unique product skus, %d duplicates, %d discard values\n",
		a.store.UniqueCount(),
		a.store.DuplicatedCount(),
		a.handler.countInvalid,
	)
}

func (a *Application) Shutdown() {
	close(a.sigint)
	close(a.skus)
}

func listenSigInt(c chan os.Signal, cancel func()) {
	signal.Notify(c, os.Interrupt)
	go func() {
		_, ok := <-c
		if ok {
			fmt.Println("Got SIGINT, stopping server gracefully...")
		}
		cancel()
	}()
}
