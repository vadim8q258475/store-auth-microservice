package main

import (
	"fmt"

	"github.com/go-chi/jwtauth"
	"github.com/vadim8q258475/store-auth-microservice/app"
	"github.com/vadim8q258475/store-auth-microservice/config"
	grpcService "github.com/vadim8q258475/store-auth-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-auth-microservice/iternal/interceptor"
	userpbv1 "github.com/vadim8q258475/store-user-microservice/gen/v1"

	authService "github.com/vadim8q258475/store-auth-microservice/iternal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO
// add cacher

func main() {
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// interceptor
	intterceptor := interceptor.NewInterceptor(logger)

	// load config
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.String())

	// auth token
	authToken := jwtauth.New("HS256", []byte(cfg.SecretKey), nil)

	// user grpc client
	conn, err := grpc.NewClient(cfg.UserHost+":"+cfg.UserPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	fmt.Println(conn.GetState())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := userpbv1.NewUserServiceClient(conn)

	// service
	service := authService.NewService(client, authToken)

	// grpc service
	grpcService := grpcService.NewGrpcService(service)

	// grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			intterceptor.UnaryServerInterceptor,
		),
	)

	// app
	app := app.NewApp(grpcService, server, logger, cfg)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
