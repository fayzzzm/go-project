package service

import (
	"context"
	"testing"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements service.UserRepository for testing.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error) {
	args := m.Called(ctx, limit, offset, teamID, tenantID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetForLogin(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestUserService_Create(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewUserService(mockRepo)

	user := &domain.User{
		Email: "test@example.com",
		Name:  strPtr("Test User"),
	}

	// Setup expectations
	mockRepo.On("Create", mock.Anything, user).Return(nil)

	// Execute
	err := svc.Create(context.Background(), user)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}
