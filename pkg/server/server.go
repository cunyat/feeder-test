package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
)

// Server manages incoming connections
type Server struct {
	addr string
	// maxConn defines the maximum number of concurrent connections
	maxConn int
	out     chan string
}

// New returns a new instance of a Server
func New(addr string, maxConn int, out chan string) *Server {
	return &Server{addr: addr, maxConn: maxConn, out: out}
}

func (s *Server) Start(ctx context.Context) error {

	errs := make(chan error, 2)

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("could not listen in %s: %w", s.addr, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < s.maxConn; i++ {
		go listen(ctx, ln, cancel, s.out, errs)
	}

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return nil
	}
}

func listen(ctx context.Context, ln net.Listener, cancel func(), out chan string, errs chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				errs <- fmt.Errorf("error accepting a new connection: %s", err.Error())
			}

			for {
				msg, err := bufio.NewReader(conn).ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					errs <- fmt.Errorf("could not read incomming message: %s", err.Error())
				}

				out <- msg
			}
		}
	}
}
