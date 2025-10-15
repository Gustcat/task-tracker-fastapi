package converter

import (
	"auth/internal/model"
	desc "auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) (int64, *desc.UserInfo, *timestamppb.Timestamp, *timestamppb.Timestamp) {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return user.ID, ToUserInfoFromService(user.Info), updatedAt, timestamppb.New(user.CreatedAt)
}

func ToUserInfoFromService(userinfo model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:  userinfo.Name,
		Email: userinfo.Email,
		Role:  desc.Role(userinfo.Role),
	}
}

func ToUserInfoFromDesc(userinfo *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  userinfo.Name,
		Email: userinfo.Email,
		Role:  int32(userinfo.Role),
	}
}
