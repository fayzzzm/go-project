package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthUseCase interface {
	Login(ctx context.Context, input usecase.LoginInput) (*usecase.LoginOutput, error)
}

type AuthController struct {
	uc AuthUseCase
}

func NewAuthController(r gin.IRouter, uc AuthUseCase) {
	c := &AuthController{uc: uc}

	auth := r.Group("/auth")
	{
		auth.POST("/login", middleware.BindJSON[usecase.LoginInput](), utils.Handle(c.Login, http.StatusOK))
	}
}

func (c *AuthController) Login(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.LoginInput](ctx)
	return c.uc.Login(ctx.Request.Context(), input)
}
