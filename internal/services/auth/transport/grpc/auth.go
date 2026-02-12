package grpc

import (
	"context"
	"fmt"

	authpb "github.com/anton-chornobai/stock-protos/auth/gen"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (err error)
	Logout()
	SignUp(
		ctx context.Context,
		email string,
		password string,
	) (userId string, err error)
	UpdateProfile()
	DeleteAccount()
}

type AuthService struct {
	authpb.UnimplementedAuthServer
	auth Auth
}

func (a *AuthService) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	err := a.auth.Login(ctx, req.Email, req.Password)

	if err != nil {
		return nil, fmt.Errorf("couldnt generate token")
	}

	return  &authpb.LoginResponse{}, nil
}