package util

import (
	"context"
	"slido-clone-backend/internal/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// TokenUtil provides methods for creating and parsing JWT tokens.
type TokenUtil struct {
	SecretKey string
	Redis     *redis.Client
}

// NewTokenUtil creates a new instance of TokenUtil with the provided secret key.
func NewTokenUtil(secretKey string, redisClient *redis.Client) *TokenUtil {
	return &TokenUtil{
		SecretKey: secretKey,
		Redis:     redisClient,
	}
}

// CreateToken generates a JWT token for the given Auth model.
func (t *TokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, error) {
	// create structure of the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  auth.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).UnixMilli(),
	})

	// sign the token with the secret key
	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	// store token in redis with expiration same as token expiration
	_, err = t.Redis.Set(ctx, jwtToken, auth.ID, time.Hour*24*30).Result()
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// ParseToken validates and parses the JWT token string and returns the Auth model.
func (t *TokenUtil) ParseToken(ctx context.Context, tokenString string) (*model.Auth, error) {
	// parse the token to normal structure of jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// extract claims from the token
	claims := token.Claims.(jwt.MapClaims)
	id := claims["id"].(string)
	exp := claims["exp"].(float64)

	// check if token is expired
	if int64(exp) < time.Now().UnixMilli() {
		return nil, fiber.ErrUnauthorized
	}

	// check if token exists in redis
	result, err := t.Redis.Exists(ctx, tokenString).Result()
	if err != nil {
		return nil, err
	}

	// if token not found in redis, return unauthorized
	if result == 0 {
		return nil, fiber.ErrUnauthorized
	}

	auth := &model.Auth{
		ID: id,
	}

	return auth, nil
}
