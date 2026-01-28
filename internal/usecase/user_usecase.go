package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
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
}

type UpdateUserInput struct {
	ID           string                 `json:"-"`
	Email        *string                `json:"email"`
	Name         *string                `json:"name"`
	Address      *string                `json:"address"`
	Phone        *string                `json:"phone"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

type UserIDInput struct {
	ID string `json:"-"`
}

type ListUserInput struct {
	Pagination domain.Pagination
	TeamID     string
}

type UserOutput struct {
	ID           string                 `json:"id"`
	Email        string                 `json:"email"`
	Name         string                 `json:"name"`
	Address      string                 `json:"address"`
	Role         string                 `json:"role"`
	TenantID     string                 `json:"tenant_id"`
	AppMetadata  map[string]interface{} `json:"app_metadata"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
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

	userID := domain.UserFromContext(ctx)

	user := &domain.User{
		Email:        input.Email,
		Name:         &input.Name,
		Address:      &input.Address,
		Phone:        &input.Phone,
		Password:     string(hashedBytes),
		AppMetadata:  input.AppMetadata,
		UserMetadata: input.UserMetadata,
		CreatedBy:    utils.StringPtrOrNil(userID),
		UpdatedBy:    utils.StringPtrOrNil(userID),
	}

	// Extract role from user_metadata if present
	if r, ok := input.UserMetadata["role"].(string); ok {
		user.Role = r
	} else {
		user.Role = domain.RoleUser
	}

	if err := uc.svc.Create(ctx, user); err != nil {
		return nil, err
	}

	return toUserOutput(user), nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, input UserIDInput) (*UserOutput, error) {
	user, err := uc.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return toUserOutput(user), nil
}

func (uc *UserUseCase) List(ctx context.Context, input ListUserInput) ([]UserOutput, error) {
	tenantID := domain.TenantFromContext(ctx)
	input.Pagination.Normalize()
	users, err := uc.svc.List(ctx, input.Pagination.Limit, input.Pagination.Offset, input.TeamID, tenantID)
	if err != nil {
		return nil, err
	}

	output := make([]UserOutput, len(users))
	for i, u := range users {
		output[i] = *toUserOutput(&u)
	}
	return output, nil
}

func (uc *UserUseCase) Update(ctx context.Context, input UpdateUserInput) (*UserOutput, error) {
	user, err := uc.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if input.AppMetadata != nil {
		user.AppMetadata = input.AppMetadata
	}
	if input.UserMetadata != nil {
		user.UserMetadata = input.UserMetadata
		// Update role if changed in user_metadata
		if r, ok := input.UserMetadata["role"].(string); ok {
			user.Role = r
		}
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

	if userID := domain.UserFromContext(ctx); userID != "" {
		user.UpdatedBy = utils.StringPtrOrNil(userID)
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
		ID:           u.ID,
		Email:        u.Email,
		Name:         name,
		Address:      address,
		Role:         u.Role,
		TenantID:     utils.StringValue(u.TenantID),
		AppMetadata:  u.AppMetadata,
		UserMetadata: u.UserMetadata,
	}
}

func (uc *UserUseCase) Delete(ctx context.Context, input UserIDInput) error {
	return uc.svc.Delete(ctx, input.ID)
}
