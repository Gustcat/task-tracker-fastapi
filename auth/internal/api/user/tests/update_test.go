package tests

import (
	"auth/internal/api/user"
	"auth/internal/service"
	servicemocks "auth/internal/service/mocks"
	desc "auth/pkg/user_v1"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.UpdateRequest
	}

	var (
		ctx = context.Background()

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		req = &desc.UpdateRequest{
			Id:    id,
			Name:  name,
			Email: email,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		expected        *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expected: &emptypb.Empty{},
			err:      nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, name, email).Return(nil)
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
				mock.UpdateMock.Expect(ctx, id, name, email).Return(serviceErr)
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
			resp, err := api.Update(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.expected, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
