package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/gin-gonic/gin"
)

type DeviceUseCase interface {
	Create(ctx context.Context, input usecase.CreateDeviceInput) (*usecase.DeviceOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.DeviceOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.DeviceOutput, error)
}

type DeviceController struct {
	uc DeviceUseCase
}

func NewDeviceController(r gin.IRouter, uc DeviceUseCase) {
	c := &DeviceController{uc: uc}

	devices := r.Group("/devices")
	{
		devices.POST("", middleware.BindJSON[usecase.CreateDeviceInput](), Handle(c.Create, http.StatusCreated))
		devices.GET("/:id", Handle(c.GetByID, http.StatusOK))
		devices.GET("", Handle(c.List, http.StatusOK))
	}
}

func (c *DeviceController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateDeviceInput](ctx)
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
