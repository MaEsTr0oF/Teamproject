package main

import (
	"TeamProject/api/auth/internal/app"
	"TeamProject/api/auth/internal/config"
	"TeamProject/api/auth/internal/storage/postgres"
	. "TeamProject/api/lib/constants"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustConfig()

	log := setepLogger(cfg.Env)

	storage, err := postgres.New(cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Host)
	if err != nil {
		panic(err)
	}

	application := app.New(log, cfg.GRPC.Port, cfg.GRPC.Timeout, storage)
	go func() {
		application.GRPCServer.MustLoad()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}

func setepLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case LocalEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case DevEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case ProdEnv:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return log
}
