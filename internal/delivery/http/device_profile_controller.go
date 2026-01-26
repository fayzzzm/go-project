package http

import (
	"context"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/reply"
	"github.com/gin-gonic/gin"
)

type DeviceProfileUseCase interface {
	Create(ctx context.Context, input usecase.CreateDeviceProfileInput) (*usecase.DeviceProfileOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.DeviceProfileOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.DeviceProfileOutput, error)
	Update(ctx context.Context, id string, input usecase.UpdateDeviceProfileInput) (*usecase.DeviceProfileOutput, error)
	Delete(ctx context.Context, id string) error
}

type DeviceProfileController struct {
	uc DeviceProfileUseCase
}

func NewDeviceProfileController(r gin.IRouter, uc DeviceProfileUseCase) {
	c := &DeviceProfileController{uc: uc}

	dps := r.Group("/device_profiles")
	{
		dps.POST("", middleware.BindJSON[usecase.CreateDeviceProfileInput](), c.Create)
		dps.GET("", c.List)
		dps.GET("/:id", c.GetByID)
		dps.PUT("/:id", middleware.BindJSON[usecase.UpdateDeviceProfileInput](), c.Update)
		dps.DELETE("/:id", c.Delete)
	}
}

func (c *DeviceProfileController) Create(ctx *gin.Context) {
	input := middleware.GetBody[usecase.CreateDeviceProfileInput](ctx)

	output, err := c.uc.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.Created(ctx, output)
}

func (c *DeviceProfileController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	output, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *DeviceProfileController) List(ctx *gin.Context) {
	limit := 100
	offset := 0

	output, err := c.uc.List(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *DeviceProfileController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateDeviceProfileInput](ctx)

	output, err := c.uc.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.OK(ctx, output)
}

func (c *DeviceProfileController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.Deleted(ctx, gin.H{"id": id, "deleted": true})
}
