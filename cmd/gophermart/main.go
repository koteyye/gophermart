package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/sergeizaitcev/gophermart/internal/gophermart"
)

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	cmd := gophermart.NewCommand()

	err := cmd.Parse(os.Args[1:])
	if err != nil {
		cmd.Usage()
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	return cmd.Run(ctx)
}
