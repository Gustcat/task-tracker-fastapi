package user

import (
	"auth/internal/model"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serv) Get(ctx context.Context, id int64) (int64, *model.UserInfo, *timestamppb.Timestamp, *timestamppb.Timestamp, error) {
	id, userinfo, createdAt, updatedAtTime, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return 0, nil, &timestamppb.Timestamp{}, &timestamppb.Timestamp{}, err
	}
	var updatedAt *timestamppb.Timestamp
	if updatedAtTime.Valid {
		updatedAt = timestamppb.New(updatedAtTime.Time)
	}

	return id, userinfo, timestamppb.New(createdAt), updatedAt, nil
}
