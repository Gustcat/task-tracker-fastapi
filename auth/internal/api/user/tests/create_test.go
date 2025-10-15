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
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()

		id  = gofakeit.Int64()
		pwd = gofakeit.Password(true, true, true, true, true, 10)

		name  = gofakeit.Name()
		email = gofakeit.Email()
		role  = int32(gofakeit.Number(0, 2))

		userinfo = &model.UserInfo{
			Name:  name,
			Email: email,
			Role:  role,
		}

		req = &desc.CreateRequest{
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  desc.Role(role),
			},
			Password:        pwd,
			PasswordConfirm: pwd,
		}

		res = &desc.CreateResponse{
			Id: id,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		expected        *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				req: req,
				ctx: ctx,
			},
			expected: res,
			err:      nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userinfo, pwd).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				req: req,
				ctx: ctx,
			},
			expected: nil,
			err:      serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, userinfo, pwd).Return(0, serviceErr)
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

			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}
