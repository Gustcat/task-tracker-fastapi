package tests

import (
	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/service"
	servicemocks "auth/internal/service/mocks"
	desc "auth/pkg/user_v1"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()

		id        = gofakeit.Int64()
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = int32(gofakeit.Number(0, 2))

		userinfo = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		req = &desc.GetRequest{
			Id: id,
		}

		res = &desc.GetResponse{
			Id: id,
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		expected        *desc.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expected: res,
			err:      nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(id, userinfo, timestamppb.New(createdAt), timestamppb.New(updatedAt), nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expected: nil,
			err:      serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(0, nil, &timestamppb.Timestamp{}, &timestamppb.Timestamp{}, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)
			resp, err := api.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}
