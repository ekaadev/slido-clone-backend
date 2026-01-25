package converter

import (
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
)

// UserToResponse convert entity User to model UserResponse
func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

// UserToAuthResponse convert entity User and token to model AuthResponse
func UserToAuthResponse(user *entity.User, token string) *model.AuthResponse {
	return &model.AuthResponse{
		User:  *UserToResponse(user),
		Token: token,
	}
}
