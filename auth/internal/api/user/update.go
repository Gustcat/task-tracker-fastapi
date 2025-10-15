package user

import (
	desc "auth/pkg/user_v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {

	err := i.userService.Update(ctx, req.GetId(), req.GetName(), req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
