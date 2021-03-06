package server

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)

// Server manages incoming connections
type Server struct {
	// addr defines the address to listen for incoing connections
	addr string
	// maxConn defines the maximum number of concurrent connections
	maxConn int
	// out is the channel where we will send received messages
	out chan string
}

// New returns a new instance of a Server
func New(addr string, maxConn int, out chan string) *Server {
	return &Server{addr: addr, maxConn: maxConn, out: out}
}

// Start initializes listeners and waits for new connections
func (s *Server) Start(ctx context.Context) error {
	// errs chan will allow to child routines send an encountered error back to parent.
	errs := make(chan error, 2)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("could not listen in %s: %w", s.addr, err)
	}

	for i := 0; i < s.maxConn; i++ {
		go listen(ctx, ln, s.out, errs)
	}

	// Wait until listeners returns error or context finished
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return nil
	}
}

func listen(ctx context.Context, ln net.Listener, out chan string, errs chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				errs <- fmt.Errorf("error accepting a new connection: %s", err.Error())
				return
			}

			msg, err := ioutil.ReadAll(conn)
			if err != nil {
				if err == io.EOF {
					continue
				}
				errs <- fmt.Errorf("could not read incomming message: %s", err.Error())
			}

			out <- string(msg)

			_ = conn.Close()
		}
	}
}
