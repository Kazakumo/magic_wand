package main

import (
	"log/slog"
	"os"

	"github.com/Kazakumo/magic_wand/internal/cli"
	"github.com/Kazakumo/magic_wand/internal/config"
	"github.com/Kazakumo/magic_wand/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	logger.Init(&cfg.Log)

	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
