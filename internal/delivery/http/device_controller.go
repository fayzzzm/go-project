package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type DeviceUseCase interface {
	Create(ctx context.Context, input usecase.CreateDeviceInput) (*usecase.DeviceOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.DeviceOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.DeviceOutput, error)
	Update(ctx context.Context, id string, input usecase.CreateDeviceInput) (*usecase.DeviceOutput, error)
	Delete(ctx context.Context, id, tenantID string) error
}

type DeviceController struct {
	uc DeviceUseCase
}

func NewDeviceController(r gin.IRouter, uc DeviceUseCase) {
	c := &DeviceController{uc: uc}

	devices := r.Group("/devices")
	devices.Use(middleware.AuthMiddleware())
	{
		devices.POST("", middleware.BindJSON[usecase.CreateDeviceInput](), utils.Handle(c.Create, http.StatusCreated))
		devices.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		devices.GET("", utils.Handle(c.List, http.StatusOK))
		devices.PUT("/:id", middleware.BindJSON[usecase.CreateDeviceInput](), utils.Handle(c.Update, http.StatusOK))
		devices.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *DeviceController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateDeviceInput](ctx)
	input.TenantID = middleware.GetTenantID(ctx)
	return c.uc.Create(ctx.Request.Context(), input)
}

func (c *DeviceController) GetByID(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return c.uc.GetByID(ctx.Request.Context(), id)
}

func (c *DeviceController) List(ctx *gin.Context) (any, error) {
	limit, offset := middleware.GetPagination(ctx)
	return c.uc.List(ctx.Request.Context(), limit, offset)
}

func (c *DeviceController) Update(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.CreateDeviceInput](ctx)
	input.TenantID = middleware.GetTenantID(ctx)
	return c.uc.Update(ctx.Request.Context(), id, input)
}

func (c *DeviceController) Delete(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	tenantID := middleware.GetTenantID(ctx)
	err := c.uc.Delete(ctx.Request.Context(), id, tenantID)
	return nil, err
}
