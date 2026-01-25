package model

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Role     string `json:"role" validate:"required,oneof=presenter admin"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type AnonymousUserRequest struct {
	RoomCode    string `json:"room_code" validate:"required,len=6,alphanum"`
	DisplayName string `json:"display_name" validate:"required,min=3,max=30"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}
