package tests

import (
	"auth/internal/model"
	"auth/internal/repository"
	repomocks "auth/internal/repository/mocks"
	userserv "auth/internal/service/user"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx      context.Context
		userinfo *model.UserInfo
		pwd      string
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

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name               string
		args               args
		expected           int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:      ctx,
				userinfo: userinfo,
				pwd:      pwd,
			},
			expected: id,
			err:      nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, userinfo, pwd).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:      ctx,
				userinfo: userinfo,
				pwd:      pwd,
			},
			expected: 0,
			err:      serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
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
			userRepositoryMock := tt.userRepositoryMock(mc)
			service := userserv.NewMockService(userRepositoryMock)
			respId, err := service.Create(tt.args.ctx, tt.args.userinfo, tt.args.pwd)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, respId)
		})
	}
}
