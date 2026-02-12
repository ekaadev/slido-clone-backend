package util

import (
	"context"
	"fmt"
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
	now := time.Now()
	expiryDuration := time.Hour * 24 * 30 // 30 days
	// set expiration time for the token
	auth.ExpiresAt = jwt.NewNumericDate(now.Add(expiryDuration))
	auth.IssuedAt = jwt.NewNumericDate(now)

	// create structure of the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth)

	// sign the token with the secret key
	jwtToken, err := token.SignedString([]byte(t.SecretKey))
	if err != nil {
		return "", err
	}

	// store token in redis with expiration same as token expiration
	_, err = t.Redis.Set(ctx, jwtToken, "valid", expiryDuration).Result()
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// ParseToken validates and parses the JWT token string and returns the Auth model.
func (t *TokenUtil) ParseToken(ctx context.Context, tokenString string) (*model.Auth, error) {
	// parse the token to normal structure of jwt
	token, err := jwt.ParseWithClaims(tokenString, &model.Auth{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// extract claims from the token
	claims, ok := token.Claims.(*model.Auth)
	if !ok || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	// check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
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

	return claims, nil
}

// InvalidateToken removes the token from Redis to effectively log out the user
func (t *TokenUtil) InvalidateToken(ctx context.Context, tokenString string) error {
	// delete token from redis
	_, err := t.Redis.Del(ctx, tokenString).Result()
	if err != nil {
		return err
	}
	return nil
}
