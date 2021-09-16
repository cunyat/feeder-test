package feeder

import (
	"context"
	"fmt"
	"log"
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
	store *store.Store
	// manager that handles incoming messages
	manager *Manager
	// writter writes results in a file
	writter *LogWritter

	skus   chan string
	sigint chan os.Signal
}

// New generates a new instance of Application
func New(addr string, maxConn int, ttl time.Duration) *Application {
	writter, err := NewLogWritter()
	if err != nil {
		log.Fatalf("could not open log file: %s", err)
	}

	// Initialize store and subscribe logger
	store := store.New()
	store.Subscribe(writter.Write)

	manager := NewManager(store, ValidateSKU)

	a := &Application{
		ttl:     ttl,
		store:   store,
		manager: manager,
		writter: writter,
		skus:    make(chan string, 10),
		sigint:  make(chan os.Signal, 1),
	}

	a.srv = server.New(addr, maxConn, a.skus)

	return a
}

// Start runs the app
func (a *Application) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), a.ttl)
	listenSigInt(a.sigint, cancel)

	go a.handleMessages(a.skus, a.manager, cancel)

	// run the server, it will return when encounters an error or the context is closed
	if err := a.srv.Start(ctx); err != nil {
		fmt.Printf("Got error from server: %s", err)
	}

	a.Shutdown()

	fmt.Printf(
		"Received %d unique product skus, %d duplicates, %d discard values\n",
		a.store.UniqueCount(),
		a.store.DuplicatedCount(),
		a.manager.countInvalid,
	)
}

// Shutdown closes open resources and channels
func (a *Application) Shutdown() {
	err := a.writter.Close()
	if err != nil {
		fmt.Printf("error closing log file: %s", err)
	}

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

// Handle messages listens for new messages in channel and sends them to manager
// Also listens for terminate sequence for shutting down the application
func (a *Application) handleMessages(skus chan string, manager *Manager, cancel context.CancelFunc) {
	for sku := range skus {
		if IsTerminateSequence(sku) {
			fmt.Println("Got terminate sequece, stopping server gracefully...")
			cancel()
			return
		}

		manager.HandleMessage(sku)
	}
}
