package grpc

import (
	"context"

	gen "github.com/vadim8q258475/store-auth-microservice/gen/v1"
	"github.com/vadim8q258475/store-auth-microservice/iternal/service"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type GrpcService struct {
	gen.UnimplementedAuthServiceServer
	service service.Service
}

func NewGrpcService(service service.Service) *GrpcService {
	return &GrpcService{
		service: service,
	}
}

func (g *GrpcService) IsTokenValid(ctx context.Context, request *gen.IsTokenValid_Request) (*gen.IsTokenValid_Response, error) {
	id, err := g.service.IsTokenValid(ctx, request.Token)
	if err != nil {
		return nil, err
	}
	return &gen.IsTokenValid_Response{Id: id}, nil
}

func (g *GrpcService) Login(ctx context.Context, request *gen.Login_Request) (*gen.Login_Response, error) {
	user, err := g.service.GetUser(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, err
	}
	token, err := g.service.GenToken(user.Id)
	if err != nil {
		return nil, err
	}
	return &gen.Login_Response{Token: token}, nil
}

func (g *GrpcService) Register(ctx context.Context, request *gen.Register_Request) (*gen.Register_Response, error) {
	_, err := g.service.GetUser(ctx, request.Email)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}
	st, ok := status.FromError(err)
	if !ok {
		return nil, err
	}
	code := st.Code()
	if code != codes.NotFound {
		return nil, err
	}
	err = g.service.Create(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &gen.Register_Response{Success: true}, nil
}
