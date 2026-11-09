package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gleb76/web-messenger/internal/config"
	"github.com/Gleb76/web-messenger/internal/server"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := server.New(cfg, logger).Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}
