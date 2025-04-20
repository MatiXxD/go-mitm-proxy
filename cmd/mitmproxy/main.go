package main

import (
	"context"
	proxyDelivery "github.com/MatiXxD/go-mitm-proxy/internal/delivery/proxy"
	requestDelivery "github.com/MatiXxD/go-mitm-proxy/internal/delivery/request"
	proxyServer "github.com/MatiXxD/go-mitm-proxy/internal/proxy"
	proxyRepository "github.com/MatiXxD/go-mitm-proxy/internal/repository/proxy"
	requestRepository "github.com/MatiXxD/go-mitm-proxy/internal/repository/request"
	requestUsecase "github.com/MatiXxD/go-mitm-proxy/internal/usecase/request"
	"github.com/MatiXxD/go-mitm-proxy/internal/webapi"
	"github.com/MatiXxD/go-mitm-proxy/pkg/db/mongodb"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
	"github.com/MatiXxD/go-mitm-proxy/pkg/logger"
	"log"
	"os"
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

	db, err := mongodb.NewMongoDB(context.Background(), cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.NewLogger("debug")
	if err != nil {
		log.Fatal(err)
	}

	// Webapi
	rr := requestRepository.NewRequestRepository(db, logger)
	ru := requestUsecase.NewRequestUsecase(rr, logger)
	rd := requestDelivery.NewRequestDelivery(ru, logger)
	webapi := webapi.NewServer(logger, cfg)
	webapi.BindRoutes(rd)

	// Proxy
	pr := proxyRepository.NewMemProxyRepository()
	pd, err := proxyDelivery.NewProxyDelivery(pr, ru, cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	proxy := proxyServer.NewProxy(pd, cfg, logger)

	// Run servers
	errChan := make(chan error, 2)

	go func() {
		errChan <- webapi.Run()
	}()

	go func() {
		errChan <- proxy.Start()
	}()

	if err := <-errChan; err != nil {
		log.Fatal(err)
	}
}
