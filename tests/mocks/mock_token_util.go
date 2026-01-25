package mocks

import (
	"context"
	"slido-clone-backend/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockTokenUtil mock untuk TokenUtil
type MockTokenUtil struct {
	mock.Mock
}

// CreateToken mock implementation
func (m *MockTokenUtil) CreateToken(ctx context.Context, auth *model.Auth) (string, error) {
	args := m.Called(ctx, auth)
	return args.String(0), args.Error(1)
}

// ParseToken mock implementation
func (m *MockTokenUtil) ParseToken(ctx context.Context, tokenString string) (*model.Auth, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Auth), args.Error(1)
}
