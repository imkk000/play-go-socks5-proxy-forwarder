package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

var addr = ":9001"

func main() {
	inputChains := flag.String("chains", "", "set proxy chains")
	enabledLog := flag.Bool("log", false, "enable log")
	flag.Parse()

	var defaultDial proxy.Dialer = &net.Dialer{}
	if inputChains != nil && len(*inputChains) > 0 {
		chains := strings.Split(*inputChains, ",")
		previousDialer, _ := proxy.SOCKS5("tcp", chains[0], nil, proxy.Direct)
		if *enabledLog {
			fmt.Printf("[%s] chain 1: %s\n", time.Now().Format(time.RFC3339Nano), chains[0])
		}
		for i, next := range chains[1:] {
			previousDialer, _ = proxy.SOCKS5("tcp", next, nil, previousDialer)
			if *enabledLog {
				fmt.Printf("[%s] chain %d: %s\n", time.Now().Format(time.RFC3339Nano), i+2, next)
			}
		}
		defaultDial = previousDialer
	}
	dialer := defaultDial.(proxy.ContextDialer)
	srv := socks5.NewServer(socks5.WithDialAndRequest(func(ctx context.Context, network, addr string, req *socks5.Request) (net.Conn, error) {
		if *enabledLog {
			fmt.Printf("[%s] from: %s -> %s (%b)\n", time.Now().Format(time.RFC3339Nano), req.LocalAddr.String(), addr, req.Request.Command)
		}

		return dialer.DialContext(ctx, network, addr)
	}))

	if *enabledLog {
		fmt.Printf("[%s] start on %s\n", time.Now().Format(time.RFC3339Nano), addr)
	}
	go srv.ListenAndServe("tcp", addr)

	done := make(chan os.Signal, 1)
	defer close(done)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	<-done
}
