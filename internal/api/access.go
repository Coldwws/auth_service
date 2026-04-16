package api

import (
	desc "authorization_service/pkg/access_v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	if req.GetEndpointAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "endpoint address is required")
	}

	err := s.AccessService().Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return &emptypb.Empty{}, nil
}
