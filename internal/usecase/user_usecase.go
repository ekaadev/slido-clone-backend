package usecase

import (
	"context"
	"slido-clone-backend/internal/entity"
	"slido-clone-backend/internal/model"
	"slido-clone-backend/internal/model/converter"
	"slido-clone-backend/internal/repository"
	"slido-clone-backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	UserRepository        *repository.UserRepository
	ParticipantRepository *repository.ParticipantRepository
	RoomRepository        *repository.RoomRepository
	TokenUtil             *util.TokenUtil
}

// NewUserUseCase create new instance of UserUseCase
func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, participantRepository *repository.ParticipantRepository, roomRepository *repository.RoomRepository, tokenUtil *util.TokenUtil) *UserUseCase {
	return &UserUseCase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		UserRepository:        userRepository,
		ParticipantRepository: participantRepository,
		RoomRepository:        roomRepository,
		TokenUtil:             tokenUtil,
	}
}

// Create usecase untuk membuat user baru
func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.AuthResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to create user
	existingUser, err := c.UserRepository.FindByEmailOrUsername(tx, request.Email, request.Username)
	if err != nil {
		c.Log.Errorf("Failed to find user by email or username: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if existingUser != nil {
		if existingUser.Email == request.Email {
			c.Log.Warnf("Email already in use: %s", request.Email)
			return nil, fiber.NewError(fiber.StatusConflict, "Email already in use")
		}

		if existingUser.Username == request.Username {
			c.Log.Warnf("Username already in use: %s", request.Username)
			return nil, fiber.NewError(fiber.StatusConflict, "Username already in use")
		}
	}

	// generate password hash
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to generate password: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// assign to entity
	user := &entity.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: string(password),
		Role:         request.Role,
	}

	// create user in repository
	if err = c.UserRepository.Create(tx, user); err != nil {
		c.Log.Errorf("Failed to create user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// create a token jwt
	token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		UserID:      &user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		IsAnonymous: false,
	})
	if err != nil {
		c.Log.Warnf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return and convert to user response
	return converter.UserToAuthResponse(user, token), nil
}

// Login usecase untuk melakukan login user
func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.AuthResponse, error) {
	// begin db transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to find user by username
	existingUser, err := c.UserRepository.FindByUsername(tx, request.Username)
	if err != nil {
		c.Log.Errorf("Failed to find user by username: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if existingUser == nil {
		c.Log.Warnf("User not found: %s", request.Username)
		return nil, fiber.ErrUnauthorized
	}

	// compare password hash
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(request.Password))
	if err != nil {
		c.Log.Warnf("Invalid password for user: %s", request.Username)
		return nil, fiber.ErrUnauthorized
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// create a token jwt
	token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		UserID:      &existingUser.ID,
		Email:       existingUser.Email,
		Role:        existingUser.Role,
		Username:    existingUser.Username,
		IsAnonymous: false,
	})
	if err != nil {
		c.Log.Warnf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return and convert to user response
	return converter.UserToAuthResponse(existingUser, token), nil
}

// Anon usecase untuk membuat user anonymous
// dapat membuat user anonymous ketika room code valid
func (c *UserUseCase) Anon(ctx context.Context, request *model.AnonymousUserRequest) (*model.JoinRoomResponse, error) {
	// begin transaction
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// logic to create anonymous user
	// check if room code exists
	roomExisting, err := c.RoomRepository.FindByRoomCode(tx, request.RoomCode)
	if err != nil {
		c.Log.Errorf("Failed to find room by room code: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if roomExisting == nil {
		c.Log.Warnf("Room not found with code: %s", request.RoomCode)
		return nil, fiber.ErrNotFound
	}

	if roomExisting.Status == "closed" {
		c.Log.Warnf("Room is closed with code: %s", request.RoomCode)
		return nil, fiber.ErrBadRequest
	}

	anon := true

	// create participant entity
	participant := &entity.Participant{
		RoomID:      roomExisting.ID,
		DisplayName: request.DisplayName,
		IsAnonymous: &anon,
	}

	// create participant in repository
	if err = c.ParticipantRepository.Create(tx, participant); err != nil {
		c.Log.Errorf("Failed to create participant: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err = tx.Commit().Error; err != nil {
		c.Log.Errorf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Anonymous users are never room owners (always audience)
	isRoomOwner := false

	// create a token jwt for anonymous user
	token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		ParticipantID: &participant.ID,
		RoomID:        &roomExisting.ID,
		DisplayName:   participant.DisplayName,
		Role:          "anonymous",
		IsAnonymous:   *participant.IsAnonymous,
		IsRoomOwner:   isRoomOwner,
	})
	if err != nil {
		c.Log.Errorf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return and convert to join room response with role (always audience for anonymous)
	return converter.ParticipantToJoinRoomResponseWithRole(participant, token, isRoomOwner), nil
}

// Logout usecase untuk logout user dengan invalidate token
func (c *UserUseCase) Logout(ctx context.Context, tokenString string) error {
	// invalidate token di redis
	if err := c.TokenUtil.InvalidateToken(ctx, tokenString); err != nil {
		c.Log.Errorf("Failed to invalidate token: %+v", err)
		return fiber.ErrInternalServerError
	}
	return nil
}
