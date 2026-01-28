package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserServicer interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
	GetForLogin(ctx context.Context, email string) (*domain.User, error)
}

type CreateUserInput struct {
	Email        string                 `json:"email" binding:"required,email"`
	Name         string                 `json:"name" binding:"required"`
	Address      string                 `json:"address"`
	Phone        string                 `json:"phone"`
	Password     string                 `json:"password" binding:"required,min=6"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
	TenantID     string                 `json:"tenant_id"`
}

type UpdateUserInput struct {
	Email        *string                `json:"email"`
	Name         *string                `json:"name"`
	Address      *string                `json:"address"`
	Phone        *string                `json:"phone"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

type UserOutput struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

type UserUseCase struct {
	svc UserServicer
}

func NewUserUseCase(svc UserServicer) *UserUseCase {
	return &UserUseCase{
		svc: svc,
	}
}

func (uc *UserUseCase) Create(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        input.Email,
		Name:         &input.Name,
		Address:      &input.Address,
		Phone:        &input.Phone,
		Password:     string(hashedBytes),
		AppMetadata:  input.AppMetadata,
		UserMetadata: input.UserMetadata,
		TenantID:     input.TenantID,
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

func (uc *UserUseCase) List(ctx context.Context, p domain.Pagination, teamID, tenantID string) ([]UserOutput, error) {
	p.Normalize()
	users, err := uc.svc.List(ctx, p.Limit, p.Offset, teamID, tenantID)
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
	user, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.AppMetadata != nil {
		user.AppMetadata = input.AppMetadata
	}
	if input.UserMetadata != nil {
		user.UserMetadata = input.UserMetadata
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Name != nil {
		user.Name = input.Name
	}
	if input.Address != nil {
		user.Address = input.Address
	}
	if input.Phone != nil {
		user.Phone = input.Phone
	}

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
		ID:       u.ID,
		Email:    u.Email,
		Name:     name,
		Address:  address,
		Role:     u.Role,
		TenantID: u.TenantID,
	}
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}
