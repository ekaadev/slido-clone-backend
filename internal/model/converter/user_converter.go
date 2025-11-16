package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func UserToAuthResponse(user *entity.User, token string) *model.AuthResponse {
	return &model.AuthResponse{
		User:  *UserToResponse(user),
		Token: token,
	}
}
