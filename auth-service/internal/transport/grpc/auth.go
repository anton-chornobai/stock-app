package grpchandler

import (
	authpb "github.com/anton-chornobai/stock-protos/auth/gen"
)

type Auth interface {
	ValidateToken()
	UpdateProfile()
	CheckAdmin()
}

type AuthHandler struct {
	authpb.UnimplementedAuthServer
	auth Auth
}

func NewAuthHandler(auth Auth) *AuthHandler{
	return &AuthHandler{
		auth: auth,
	}
}