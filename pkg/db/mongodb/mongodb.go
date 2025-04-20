package mongodb

import (
	"context"
	"fmt"
	"github.com/MatiXxD/go-mitm-proxy/pkg/env"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 10 * time.Second
	pingTimeout     = 5 * time.Second
	maxPoolSize     = 100
	minPoolSize     = 10
	connIdleTimeout = 30 * time.Second
)

func NewMongoDB(ctx context.Context, cfg *env.Config) (*mongo.Database, error) {
	connString := fmt.Sprintf("mongodb://%s:%s",
		cfg.MongoConfig.Host,
		cfg.MongoConfig.Port,
	)

	clientOpts := options.Client().
		ApplyURI(connString).
		SetMaxPoolSize(maxPoolSize).
		SetMinPoolSize(minPoolSize).
		SetMaxConnIdleTime(connIdleTimeout)

	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client, err := mongo.Connect(connCtx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("mongodb connect error: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, pingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongodb ping error: %w", err)
	}

	// create database if it doesn't exist
	db := client.Database(cfg.MongoConfig.Database)
	return db, nil
}
