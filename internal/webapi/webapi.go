package webapi

import (
	"context"
	"errors"
	"fmt"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo   *echo.Echo
	db     *mongo.Database
	logger *zap.Logger
	cfg    *env.Config
}

func NewServer(logger *zap.Logger, cfg *env.Config) *Server {
	return &Server{
		echo:   echo.New(),
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:              s.cfg.ServerConfig.Addr,
		Handler:           s.echo,
		ReadHeaderTimeout: s.cfg.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout:      s.cfg.ServerConfig.WriteTimeout * time.Second,
		IdleTimeout:       s.cfg.ServerConfig.IdleTimeout * time.Second,
	}

	s.logger.Info(fmt.Sprintf("server listen on %s", s.cfg.ServerConfig.Addr))
	go func() {
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
