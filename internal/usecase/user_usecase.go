package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type UserServicer interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
}

type CreateUserInput struct {
	Email        string                 `json:"email" binding:"required,email"`
	Name         string                 `json:"name" binding:"required"`
	Address      string                 `json:"address"`
	Phone        string                 `json:"phone"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

type UpdateUserInput struct {
	Email        string                 `json:"email"`
	Name         string                 `json:"name"`
	Address      string                 `json:"address"`
	Phone        string                 `json:"phone"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

type UserOutput struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type UserUseCase struct {
	svc UserServicer
}

func NewUserUseCase(svc UserServicer) *UserUseCase {
	return &UserUseCase{svc: svc}
}

func (uc *UserUseCase) Create(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	user := &domain.User{
		Email:        input.Email,
		Name:         &input.Name,
		Address:      &input.Address,
		Phone:        &input.Phone,
		AppMetadata:  input.AppMetadata,
		UserMetadata: input.UserMetadata,
	}

	if err := uc.svc.Create(ctx, user); err != nil {
		return nil, err
	}

	return toUserOutput(user), nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*UserOutput, error) {
	user, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toUserOutput(user), nil
}

func (uc *UserUseCase) List(ctx context.Context, limit, offset int) ([]UserOutput, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	users, err := uc.svc.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	output := make([]UserOutput, len(users))
	for i, u := range users {
		output[i] = *toUserOutput(&u)
	}
	return output, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id string, input UpdateUserInput) (*UserOutput, error) {
	// Direct Partial Update via Repository
	// DB handles COALESCE logic (only update fields that are present/non-nil)
	user := &domain.User{
		ID: id,
	}

	if input.Name != "" {
		user.Name = &input.Name
	}
	if input.Address != "" {
		user.Address = &input.Address
	}
	if input.Phone != "" {
		user.Phone = &input.Phone
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	// AppMetadata/UserMetadata handling... (omitted for brevity, can add if needed)

	// Update returns the updated user state from DB (RETURNING *)
	if err := uc.svc.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserOutput(user), nil
}

func toUserOutput(u *domain.User) *UserOutput {
	name := ""
	if u.Name != nil {
		name = *u.Name
	}
	address := ""
	if u.Address != nil {
		address = *u.Address
	}
	return &UserOutput{
		ID:      u.ID,
		Email:   u.Email,
		Name:    name,
		Address: address,
	}
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}
