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

func TestDelete(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) service.UserService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		req = &desc.DeleteRequest{
			Id: id,
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
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		}, {
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			expected: nil,
			err:      serviceErr,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := servicemocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
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
			resp, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.expected, resp)
			require.Equal(t, tt.err, err)
		})
	}
}
