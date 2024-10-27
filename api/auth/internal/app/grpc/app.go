package grpcApp

import (
	authgrpc "TeamProject/api/auth/internal/grpc/auth"
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	Log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authServer authgrpc.Auth,
	port int,
) *App {

	loggerOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainStreamInterceptor(
		recovery.StreamServerInterceptor(recoveryOpts...),
		logging.StreamServerInterceptor(interceptorLogger(log), loggerOpts...),
	),
	)

	authgrpc.Register(gRPCServer, authServer)

	return &App{
		Log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustLoad() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.Log.Info("auth service started", slog.String("addr", l.Addr().String()))

	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcApp.Stop"

	a.Log.With(slog.String("op", op)).Info("stopping auth service")

	a.gRPCServer.GracefulStop()
}

func interceptorLogger(log *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		log.Log(ctx, slog.Level(level), "", fields...)
	})
}
