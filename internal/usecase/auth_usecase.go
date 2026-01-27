package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServicer interface {
	GetForLogin(ctx context.Context, email string) (*domain.User, error)
}

// Note: UserServicer satisfies this, or we can just use UserServicer directly.
// To keep it simple, I will use UserServicer dependency but only call GetForLogin.

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginOutput struct {
	AccessToken string     `json:"access_token"`
	User        UserOutput `json:"user"`
}

type AuthUseCase struct {
	svc       UserServicer // Reusing UserServicer as it has GetForLogin
	jwtSecret string
}

func NewAuthUseCase(svc UserServicer) *AuthUseCase {
	return &AuthUseCase{
		svc:       svc,
		jwtSecret: "CHANGE_THIS_TO_ENV_VAR", // TODO: Load from config
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	user, err := uc.svc.GetForLogin(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       user.ID,
		"email":     user.Email,
		"role":      user.Role,
		"tenant_id": user.TenantID,
		"exp":       time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		AccessToken: tokenString,
		User:        *toUserOutput(user),
	}, nil
}
