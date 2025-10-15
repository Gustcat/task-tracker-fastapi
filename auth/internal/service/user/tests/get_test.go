package tests

import (
	"auth/internal/model"
	"auth/internal/repository"
	repomocks "auth/internal/repository/mocks"
	userserv "auth/internal/service/user"
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository

	type args struct {
		ctx context.Context
		id  int64
	}

	type serviceRes struct {
		id        int64
		userinfo  *model.UserInfo
		createdAt *timestamppb.Timestamp
		updatedAt *timestamppb.Timestamp
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

		resp = serviceRes{
			id:        id,
			userinfo:  userinfo,
			createdAt: timestamppb.New(createdAt),
			updatedAt: timestamppb.New(updatedAt),
		}

		respErr = serviceRes{
			id:        0,
			userinfo:  nil,
			createdAt: &timestamppb.Timestamp{},
			updatedAt: &timestamppb.Timestamp{},
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name               string
		args               args
		expected           serviceRes
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			expected: resp,
			err:      nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(id, userinfo, createdAt, sql.NullTime{Time: updatedAt, Valid: true}, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			expected: respErr,
			err:      serviceErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repomocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(0, nil, time.Time{}, sql.NullTime{Time: time.Time{}, Valid: false}, serviceErr)
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
			respId, respUserinfo, respCreatedAt, respUpdatedAt, err := service.Get(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, serviceRes{
				id:        respId,
				userinfo:  respUserinfo,
				createdAt: respCreatedAt,
				updatedAt: respUpdatedAt,
			})
		})
	}
}
