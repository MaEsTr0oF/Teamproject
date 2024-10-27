package app

import (
	grpcApp "TeamProject/api/auth/internal/app/grpc"
	"TeamProject/api/auth/internal/services/auth"
	"TeamProject/api/auth/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(
	log *slog.Logger,
	port int,
	ttl time.Duration,
	storage *postgres.Storage,
) *App {
	authServer := auth.New(log, storage, storage, storage, ttl)

	app := grpcApp.New(log, authServer, port)

	return &App{
		GRPCServer: app,
	}
}
