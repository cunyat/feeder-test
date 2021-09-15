package feeder

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
	"github.com/cunyat/feeder/pkg/store"
)

// Application holds all components of this app, initalizes them and starts execution
type Application struct {
	// time to live of the server
	ttl time.Duration
	// server that listens for incoming messages
	srv *server.Server
	// store to keep skus deduplicated
	store *store.DeduplicatedStore
	// manager that handles incoming messages
	manager *Manager
	// writter writes results in a file
	writter *OutputWritter

	skus   chan string
	sigint chan os.Signal
}

// New generates a new instance of Application
func New(addr string, maxConn int, ttl time.Duration, outfile string) *Application {
	a := &Application{
		ttl:     ttl,
		store:   store.New(),
		writter: NewOutputWritter(outfile),
		skus:    make(chan string, 5),
		sigint:  make(chan os.Signal, 1),
	}

	a.manager = NewManager(a.store, ValidateSKU)
	a.srv = server.New(addr, maxConn, a.skus)

	return a
}

// Start runs the app
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

			a.manager.HandleMessage(sku)
		}
	}(a.manager, a.skus, cancel)

	if err := a.srv.Start(ctx); err != nil {
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
		a.manager.countInvalid,
	)
}

// Shutdown closes open chans
func (a *Application) Shutdown() {
	close(a.sigint)
	close(a.skus)
}

// listenSigInt will cancel the context when CTRL+C is received
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
