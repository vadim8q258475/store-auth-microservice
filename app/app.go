package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vadim8q258475/store-auth-microservice/config"
	gen "github.com/vadim8q258475/store-auth-microservice/gen/v1"
	grpcService "github.com/vadim8q258475/store-auth-microservice/iternal/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	service *grpcService.GrpcService
	server  *grpc.Server
	logger  zap.Logger
	port    string
}

func NewApp(service *grpcService.GrpcService, server *grpc.Server, logger *zap.Logger, cfg config.Config) *App {
	return &App{
		service: service,
		port:    cfg.Port,
		server:  server,
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return err
	}
	gen.RegisterAuthServiceServer(a.server, a.service)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := a.server.Serve(l); err != nil {
			a.logger.Error("Server error", zap.Error(err))
		}
	}()

	<-stop

	a.logger.Info("Shutting down gRPC server...")
	a.server.GracefulStop()
	a.logger.Info("gRPC server stopped gracefully")

	return nil
}
