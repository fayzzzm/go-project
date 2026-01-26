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
		devices.POST("", middleware.BindJSON[usecase.CreateDeviceInput](), c.Create)
		devices.GET("/:id", c.GetByID)
		devices.GET("", c.List)
	}
}

func (c *DeviceController) Create(ctx *gin.Context) {
	input := middleware.GetBody[usecase.CreateDeviceInput](ctx)

	output, err := c.uc.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, output)
}

func (c *DeviceController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	output, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, output)
}

func (c *DeviceController) List(ctx *gin.Context) {
	devices, err := c.uc.List(ctx.Request.Context(), 10, 0)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, devices)
}
