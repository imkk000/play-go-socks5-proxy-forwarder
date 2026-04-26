package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand/v2"
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

func resolveDialer(chainAddr string, base proxy.Dialer) proxy.Dialer {
	if chainAddr == "direct" {
		return proxy.Direct
	}
	d, _ := proxy.SOCKS5("tcp", chainAddr, nil, base)
	return d
}

func main() {
	var multiChains [][]string
	flag.Func("chains", "set proxy chains (support multiple random chains)", func(s string) error {
		if len(s) == 0 {
			return errors.New("chains is empty")
		}
		chains := strings.Split(s, ",")
		multiChains = append(multiChains, chains)

		return nil
	})
	enabledLog := flag.Bool("log", false, "enable log")
	flag.Parse()

	var dialers []proxy.ContextDialer
	if len(multiChains) > 0 {
		for id, chains := range multiChains {
			previousDialer := resolveDialer(chains[0], proxy.Direct)
			if *enabledLog {
				fmt.Printf("[%s] chain %d - 1: %s\n", time.Now().Format(time.RFC3339Nano), id, chains[0])
			}
			for i, next := range chains[1:] {
				previousDialer = resolveDialer(next, previousDialer)
				if *enabledLog {
					fmt.Printf("[%s] chain %d - %d: %s\n", time.Now().Format(time.RFC3339Nano), id, i+2, next)
				}
			}
			dialer := previousDialer.(proxy.ContextDialer)
			dialers = append(dialers, dialer)
		}
	}
	if len(dialers) == 0 {
		dialers = append(dialers, &net.Dialer{})
	}
	l := len(dialers)
	getDialer := func() (proxy.ContextDialer, int) {
		var id int
		if l > 0 {
			id = rand.IntN(l)
		}
		return dialers[id], id
	}

	srv := socks5.NewServer(socks5.WithDialAndRequest(func(ctx context.Context, network, addr string, req *socks5.Request) (net.Conn, error) {
		dialer, id := getDialer()
		if *enabledLog {
			fmt.Printf("[%s] id: %d from: %s -> %s (%d)\n", time.Now().Format(time.RFC3339Nano), id, req.LocalAddr.String(), addr, req.Request.Command)
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
