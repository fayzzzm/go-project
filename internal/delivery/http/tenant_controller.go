package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TenantUseCase interface {
	List(ctx context.Context, p domain.Pagination) ([]usecase.TenantOutput, error)
}

type TenantController struct {
	uc TenantUseCase
}

func NewTenantController(r gin.IRouter, uc TenantUseCase) {
	c := &TenantController{uc: uc}

	tenants := r.Group("/tenants")
	{
		tenants.GET("", utils.Handle(c.List, http.StatusOK))
	}
}

func (c *TenantController) List(ctx *gin.Context) (any, error) {
	return c.uc.List(ctx.Request.Context(), middleware.GetPagination(ctx))
}
