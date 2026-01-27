package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
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
		dps.POST("", middleware.BindJSON[usecase.CreateDeviceProfileInput](), utils.Handle(c.Create, http.StatusCreated))
		dps.GET("", utils.Handle(c.List, http.StatusOK))
		dps.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		dps.PUT("/:id", middleware.BindJSON[usecase.UpdateDeviceProfileInput](), utils.Handle(c.Update, http.StatusOK))
		dps.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *DeviceProfileController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateDeviceProfileInput](ctx)
	return c.uc.Create(ctx.Request.Context(), input)
}

func (c *DeviceProfileController) GetByID(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return c.uc.GetByID(ctx.Request.Context(), id)
}

func (c *DeviceProfileController) List(ctx *gin.Context) (any, error) {
	limit, offset := middleware.GetPagination(ctx)
	return c.uc.List(ctx.Request.Context(), limit, offset)
}

func (c *DeviceProfileController) Update(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateDeviceProfileInput](ctx)
	return c.uc.Update(ctx.Request.Context(), id, input)
}

func (c *DeviceProfileController) Delete(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	return nil, err
}
