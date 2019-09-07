package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type ServerTCP struct {
	host    string
	port    string
	timeout int
	wg      sync.WaitGroup
}

func NewServerTCP(t int, h, p string) *ServerTCP {
	return &ServerTCP{
		host:    h,
		port:    p,
		timeout: t,
		wg:      sync.WaitGroup{},
	}
}

func (s *ServerTCP) ConnectAndServe() error {
	d := &net.Dialer{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.timeout)*time.Second)

	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
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
	//conn.Close()
	return nil
}

func (s *ServerTCP) ReadRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
READ:
	for {
		select {
		case <-ctx.Done():
			break READ
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN")
				break READ
			}
			text := scanner.Text()
			log.Printf("Server: %s", text)
		}
	}
	log.Printf("Finished readRoutine")
}

func (s *ServerTCP) WriteRoutine(ctx context.Context, conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
WRITE:
	for {
		select {
		case <-ctx.Done():
			break WRITE
		default:
			if !scanner.Scan() {
				break WRITE
			}
			str := scanner.Text()
			log.Printf("Client: %v\n", str)

			if _, err := conn.Write([]byte(fmt.Sprintf("%s\n", str))); err != nil {
				log.Panicf("error: %v", err)
				break WRITE
			}
		}

	}
	log.Printf("Finished writeRoutine")
}

func (s *ServerTCP) handleInput(inputChan chan string, errChan chan error) {
	// TODO: Move the scanner func here and handle the input intelligently!
}
