package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TenantUseCase interface {
	List(ctx context.Context, input usecase.ListTenantInput) ([]usecase.TenantOutput, error)
}

type TenantController struct {
	uc     TenantUseCase
	userUC UserUseCase
}

func NewTenantController(r gin.IRouter, uc TenantUseCase, userUC UserUseCase) {
	c := &TenantController{uc: uc, userUC: userUC}

	tenants := r.Group("/tenants")
	tenants.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
	{
		tenants.GET("", utils.Handle(c.List, http.StatusOK))

		users := tenants.Group("/users")
		{
			users.GET("", utils.Handle(c.ListTenantUsers, http.StatusOK))
		}
	}
}

func (c *TenantController) List(ctx *gin.Context) ([]usecase.TenantOutput, error) {
	return c.uc.List(ctx.Request.Context(), usecase.ListTenantInput{Pagination: middleware.GetPagination(ctx)})
}

func (c *TenantController) ListTenantUsers(ctx *gin.Context) ([]usecase.UserOutput, error) {
	return c.userUC.List(ctx.Request.Context(), usecase.ListUserInput{
		Pagination: middleware.GetPagination(ctx),
	})
}
