package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

// ServerTCP type
type ServerTCP struct {
	host    string
	port    string
	timeout int
	wg      sync.WaitGroup
}

// NewServerTCP returns new ServerTCP object to the caller
func NewServerTCP(t int, h, p string) *ServerTCP {
	return &ServerTCP{
		host:    h,
		port:    p,
		timeout: t,
		wg:      sync.WaitGroup{},
	}
}

// ConnectAndServe dials to the specified host and subtly handles unidirectional data flow
func (s *ServerTCP) ConnectAndServe() error {
	d := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.timeout)*time.Second)

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()

	// Handle interrupt
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	go func() {
		<-exitChan
		log.Fatal("Received SIGINT")
	}()

	s.wg.Add(1)
	go func() {
		s.ReadRoutine(ctx, conn)
		s.wg.Done()
	}()

	s.wg.Add(1)
	go func() {
		s.WriteRoutine(ctx, conn)
		s.wg.Done()
	}()

	s.wg.Wait()
	cancel()
	return nil
}

// ReadRoutine reads data received from the Network connection and logs it
func (s *ServerTCP) ReadRoutine(ctx context.Context, conn net.Conn) {
	output := make(chan string)
	err := make(chan error)
	go s.handleInput(conn, output, err)

READ:
	for {
		select {
		case <-ctx.Done():
			break READ
		case line, ok := <-output:
			if ok {
				log.Printf("Server: %s", line)
			}
		case e := <-err:
			log.Fatal(e)
		}
	}
	log.Printf("ReadRoutine has finished execution!")
}

// WriteRoutine writes data received from STDIN to the server recipient
func (s *ServerTCP) WriteRoutine(ctx context.Context, conn net.Conn) {
	// These two channels will control the main work cycle
	inputChan := make(chan string)
	errChan := make(chan error)
	go s.handleInput(os.Stdin, inputChan, errChan)

WRITE:
	for {
		select {
		case <-ctx.Done():
			break WRITE
		case in, ok := <-inputChan:
			if ok {
				if _, err := conn.Write([]byte(fmt.Sprintf("Client: %s", in))); err != nil {
					log.Fatalf("error: %v", err)
				}
			}
		case e := <-errChan:
			log.Fatalf("error: %v", e)
		}
	}
	log.Printf("WriteRoutine has finished execution!")
}

func (s *ServerTCP) handleInput(src io.Reader, inputChan chan string, errChan chan error) {
	r := bufio.NewReader(src)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			errChan <- err
			break
		}
		inputChan <- line
	}
}
