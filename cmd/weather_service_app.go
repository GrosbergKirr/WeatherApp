package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GrosbergKirr/WeatherApp/internal"
	"github.com/GrosbergKirr/WeatherApp/internal/server"
	"github.com/GrosbergKirr/WeatherApp/internal/storage"
)

func WeatherServiceApp(ctx context.Context, log *slog.Logger,
	cfg *internal.Config,
	db *storage.Storage) {
	router := server.SetRouters(log, db)
	newServer := server.NewServer(cfg, router)
	serverStopSig := make(chan os.Signal)
	signal.Notify(serverStopSig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go newServer.ServerRun(log, cfg)
	<-serverStopSig
	newServer.ServerStop(ctx, log)
}
