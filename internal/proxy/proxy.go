package proxy

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/MatiXxD/go-mitm-proxy/internal/delivery/proxy"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

type Proxy struct {
	delivery *proxy.ProxyDelivery
	cfg      *env.Config
}

func NewProxy(pd *proxy.ProxyDelivery, cfg *env.Config) *Proxy {
	return &Proxy{
		delivery: pd,
		cfg:      cfg,
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp4", p.cfg.Addr)
	if err != nil {
		return fmt.Errorf("can't listen on %s: %v", p.cfg.Addr, err)
	}
	log.Printf("proxy listen on %s", p.cfg.Addr)

	wg := &sync.WaitGroup{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Printf("Server closed: %v", err)
				break
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.delivery.Handle(conn); err != nil {
				log.Printf("Error while handle conn: %v", err)
			}
		}()
	}
	wg.Wait()
	return nil
}
