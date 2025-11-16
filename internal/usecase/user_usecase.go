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
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	TokenUtil      *util.TokenUtil
}

// NewUserUseCase create new instance of UserUseCase
func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository, tokenUtil *util.TokenUtil) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
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
		ID: user.Username,
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

	// create a token jwt
	token, err := c.TokenUtil.CreateToken(ctx, &model.Auth{
		ID: existingUser.Username,
	})
	if err != nil {
		c.Log.Warnf("Failed to create token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// return and convert to user response
	return converter.UserToAuthResponse(existingUser, token), nil
}
