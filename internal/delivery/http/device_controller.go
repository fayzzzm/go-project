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

type DeviceUseCase interface {
	Create(ctx context.Context, input usecase.CreateDeviceInput) (*usecase.DeviceOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.DeviceOutput, error)
	List(ctx context.Context, p domain.Pagination) ([]usecase.DeviceOutput, error)
	Update(ctx context.Context, id string, input usecase.CreateDeviceInput) (*usecase.DeviceOutput, error)
	Delete(ctx context.Context, id string) error
}

type DeviceController struct {
	uc DeviceUseCase
}

func NewDeviceController(r gin.IRouter, uc DeviceUseCase) {
	c := &DeviceController{uc: uc}

	devices := r.Group("/devices")
	devices.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
	{
		devices.POST("", middleware.BindJSON[usecase.CreateDeviceInput](), utils.Handle(c.Create, http.StatusCreated))
		devices.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		devices.GET("", utils.Handle(c.List, http.StatusOK))
		devices.PUT("/:id", middleware.BindJSON[usecase.CreateDeviceInput](), utils.Handle(c.Update, http.StatusOK))
		devices.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *DeviceController) Create(ctx *gin.Context) (*usecase.DeviceOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateDeviceInput](ctx))
}

func (c *DeviceController) GetByID(ctx *gin.Context) (*usecase.DeviceOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), ctx.Param("id"))
}

func (c *DeviceController) List(ctx *gin.Context) ([]usecase.DeviceOutput, error) {
	return c.uc.List(ctx.Request.Context(), middleware.GetPagination(ctx))
}

func (c *DeviceController) Update(ctx *gin.Context) (*usecase.DeviceOutput, error) {
	return c.uc.Update(ctx.Request.Context(), ctx.Param("id"), middleware.GetBody[usecase.CreateDeviceInput](ctx))
}

func (c *DeviceController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), ctx.Param("id"))
}
