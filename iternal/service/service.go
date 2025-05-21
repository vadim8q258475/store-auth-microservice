package service

import (
	"context"
	"errors"
	"strings"

	"github.com/go-chi/jwtauth"
	userpb "github.com/vadim8q258475/store-user-microservice/gen/v1"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetUser(ctx context.Context, email string) (*userpb.GetByEmail_Response, error)
	GenToken(id uint32) (string, error)
	Create(ctx context.Context, email, password string) error
	IsTokenValid(ctx context.Context, tokenString string) (uint32, error)
}

type service struct {
	client    userpb.UserServiceClient
	authToken *jwtauth.JWTAuth
}

func NewService(client userpb.UserServiceClient, authToken *jwtauth.JWTAuth) Service {
	return &service{
		client:    client,
		authToken: authToken,
	}
}

func (s *service) GetUser(ctx context.Context, email string) (*userpb.GetByEmail_Response, error) {
	request := &userpb.GetByEmail_Request{Email: email}
	return s.client.GetByEmail(ctx, request)
}

func (s *service) GenToken(id uint32) (string, error) {
	_, token, err := s.authToken.Encode(map[string]interface{}{"id": id})
	return token, err
}
func (s *service) Create(ctx context.Context, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.client.Create(ctx, &userpb.Create_Request{Email: email, Password: string(hashedPassword)})
	return err
}

func (s *service) IsTokenValid(ctx context.Context, tokenString string) (uint32, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := s.authToken.Decode(tokenString)
	if err != nil {
		return 0, err
	}

	claims, err := token.AsMap(ctx)
	if err != nil {
		return 0, err
	}

	id, ok := claims["id"].(uint32)
	if !ok {
		return 0, errors.New("id claim not found or invalid")
	}

	return id, nil
}
