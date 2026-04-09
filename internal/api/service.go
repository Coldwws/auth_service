package api

import (
	"authorization_service/internal/service"
	desc2 "authorization_service/pkg/access_v1"
	desc "authorization_service/pkg/auth_v1"
)

type Server struct {
	authService service.AuthService
	desc.UnimplementedAuthV1Server
	desc2.UnimplementedAccessV1Server
}

func NewServer(aService service.AuthService) *Server {
	return &Server{authService: aService}
}
