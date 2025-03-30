package main

import (
	"log"
	"os"

	proxyDelivery "github.com/MatiXxD/go-mitm-proxy/internal/delivery/proxy"
	proxyServer "github.com/MatiXxD/go-mitm-proxy/internal/proxy"
	proxyRepository "github.com/MatiXxD/go-mitm-proxy/internal/repository/proxy"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
)

func main() {
	envPath := "config/dev.env"
	if os.Getenv("CONFIG_PATH") != "" {
		envPath = os.Getenv("CONFIG_PATH")
	}
	cfg, err := env.NewConfig(envPath)
	if err != nil {
		log.Fatal(err)
	}

	pr := proxyRepository.NewMemProxyRepository()
	pd, err := proxyDelivery.NewProxyDelivery(pr, cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := proxyServer.NewProxy(pd, cfg)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
