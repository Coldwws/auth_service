package api

import (
	"authorization_service/internal/service"
	desc "authorization_service/pkg/auth_v1"
)

type Server struct {
	authService service.AuthService
	desc.UnimplementedAuthV1Server
}

func NewServer(aService service.AuthService) *Server {
	return &Server{authService: aService}
}
