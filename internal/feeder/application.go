package feeder

import (
	"context"
	"fmt"
	"github.com/cunyat/feeder/pkg/store"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/cunyat/feeder/pkg/server"
)

type Application struct {
	ttl     time.Duration
	srv     *server.Server
	store   *store.DeduplicatedStore
	handler *Manager

	skus   chan string
	sigint chan os.Signal
}

func New(addr string, maxConn int, store *store.DeduplicatedStore, ttl time.Duration) *Application {
	a := &Application{
		ttl:    ttl,
		store: store,
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

	go func(handler *Manager, skus chan string, cancel context.CancelFunc) {
		for sku := range skus {
			if IsTerminateSequence(sku) {
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

	fmt.Printf(
		"Received %d unique product skus, %d duplicates, %d discard values\n",
		a.store.UniqueCount(),
		a.store.DuplicatedCount(),
		a.handler.countInvalid,
	)

	file, err := os.OpenFile("out/skus.txt", os.O_RDWR | os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Could not write output file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file", err)
		}
	}(file)

	_, err = io.Copy(file, a.store.GetReader())
	if err != nil {
		fmt.Println("Error writting to output file")
	}
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
