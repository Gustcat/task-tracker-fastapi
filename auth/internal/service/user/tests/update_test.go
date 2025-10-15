package tests

import (
	"auth/internal/repository"
	repomocks "auth/internal/repository/mocks"
	userserv "auth/internal/service/user"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx   context.Context
		id    int64
		name  string
		email string
	}

	var (
		ctx = context.Background()

		id    = gofakeit.Int64()
		name  = gofakeit.Name()
		email = gofakeit.Email()

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:   ctx,
				id:    id,
				name:  name,
				email: email,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, name, email).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx:   ctx,
				id:    id,
				name:  name,
				email: email,
			},
			err: serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
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
			userRepositoryMock := tt.userRepositoryMock(mc)
			service := userserv.NewMockService(userRepositoryMock)
			err := service.Update(tt.args.ctx, tt.args.id, tt.args.name, tt.args.email)
			require.Equal(t, tt.err, err)
		})
	}
}
