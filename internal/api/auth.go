package api

import (
	"authorization_service/internal/model"
	desc "authorization_service/pkg/auth_v1"
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	refreshToken, err := s.AuthService().Login(ctx, &model.Login{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.Wrap(err, "login failed").Error())
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}

func (s *Server) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	newToken, err := s.AuthService().GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.Wrap(err, "failed to refresh token").Error())
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: newToken}, nil
}

func (s *Server) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	accessToken, err := s.AuthService().GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.Wrap(err, "failed to get access token").Error())
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}

func (s *Server) Register(ctx context.Context, req *desc.RegisterRequest) (*desc.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	role := req.GetRole()
	if role == "" {
		role = "user"
	}

	id, err := s.authService.Register(ctx, req.GetEmail(), req.GetPassword(), role)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "register failed").Error())
	}

	return &desc.RegisterResponse{Id: id}, nil
}
