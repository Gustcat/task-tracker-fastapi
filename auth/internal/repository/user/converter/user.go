package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
	"database/sql"
	"time"
)

func ToUserFromRepo(user *modelRepo.User) (int64, *model.UserInfo, time.Time, sql.NullTime) {
	return user.ID, ToUserInfoFromRepo(user.Info), user.CreatedAt, user.UpdatedAt
}

func ToUserInfoFromRepo(userinfo modelRepo.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:  userinfo.Name,
		Email: userinfo.Email,
		Role:  userinfo.Role,
	}
}
