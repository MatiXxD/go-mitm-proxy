package proxy

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net"
	"sync"

	"github.com/MatiXxD/go-mitm-proxy/internal/delivery/proxy"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

type Proxy struct {
	delivery *proxy.ProxyDelivery
	cfg      *env.Config
	logger   *zap.Logger
}

func NewProxy(pd *proxy.ProxyDelivery, cfg *env.Config, logger *zap.Logger) *Proxy {
	return &Proxy{
		delivery: pd,
		cfg:      cfg,
		logger:   logger,
	}
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp4", p.cfg.ProxyConfig.Addr)
	if err != nil {
		return fmt.Errorf("can't listen on %s: %v", p.cfg.ProxyConfig.Addr, err)
	}
	p.logger.Info(fmt.Sprintf("proxy listen on %s", p.cfg.ProxyConfig.Addr))

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
