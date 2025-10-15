package user

import (
	"auth/internal/converter"
	desc "auth/pkg/user_v1"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "password and password_confirm do not match")
	}

	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req.GetInfo()), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
