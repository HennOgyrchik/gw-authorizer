package main

import (
	"context"
	"flag"
	"gw-authorizer/internal/app"
	"gw-authorizer/internal/config"
	"gw-authorizer/internal/grpcServer"
	"gw-authorizer/internal/storages/postgres"
	"gw-authorizer/pkg/logs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	logger := logs.New(os.Stdout)

	confPath := flag.String("c", "config.env", "path to configuration")
	migrationPath := flag.String("m", "migrations", "path to migration DB files")
	flag.Parse()

	if err := config.LoadConfig(*confPath); err != nil {
		logger.Err("read configuration", err)
		return
	}

	cfg := config.New()

	dbUrl, err := cfg.Postgres.ConnectionURL()
	if err != nil {
		logger.Err("read db url", err)
		return
	}

	db := postgres.New()

	if err = db.Start(ctx, dbUrl, time.Duration(cfg.Postgres.ConnTimeout)*time.Second, *migrationPath); err != nil {
		logger.Err("connection db", err)
		return
	}
	defer db.Stop()

	srv := grpcServer.New(cfg.GRPC.ConnectionURL(), time.Duration(cfg.GRPC.Timeout)*time.Second, app.New(logger, db, cfg.App.CostEncoding, cfg.App.SecretKey))

	go func() {
		<-ctx.Done()
		logger.Info("Closing", logs.Attr{Key: "code", Value: "0"})
		srv.Stop()
	}()

	err = srv.Run()
	if err != nil {
		logger.Err("run gRPC", err)
		return
	}

}
