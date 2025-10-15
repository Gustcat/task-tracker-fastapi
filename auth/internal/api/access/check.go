package access

import (
	descAccess "auth/pkg/access_v1"
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
)

func (i *Implementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], i.tokenConfig.AuthPrefix()) {
		return nil, errors.New("authorization header is not bearer token")
	}

	accessToken := strings.TrimPrefix(authHeader[0], i.tokenConfig.AuthPrefix())

	err := i.accessService.Check(ctx, req.GetEndpointAddress(), accessToken)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
