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
	GetByID(ctx context.Context, input usecase.DeviceProfileIDInput) (*usecase.DeviceProfileOutput, error)
	List(ctx context.Context, input usecase.ListDeviceProfileInput) ([]usecase.DeviceProfileOutput, error)
	Update(ctx context.Context, input usecase.UpdateDeviceProfileInput) (*usecase.DeviceProfileOutput, error)
	Delete(ctx context.Context, input usecase.DeviceProfileIDInput) error
}

type DeviceProfileController struct {
	uc DeviceProfileUseCase
}

func NewDeviceProfileController(r gin.IRouter, uc DeviceProfileUseCase) {
	c := &DeviceProfileController{uc: uc}

	dps := r.Group("/device_profiles")
	dps.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
	{
		dps.POST("", middleware.BindJSON[usecase.CreateDeviceProfileInput](), utils.Handle(c.Create, http.StatusCreated))
		dps.GET("", utils.Handle(c.List, http.StatusOK))
		dps.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		dps.PUT("/:id", middleware.BindJSON[usecase.UpdateDeviceProfileInput](), utils.Handle(c.Update, http.StatusOK))
		dps.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *DeviceProfileController) Create(ctx *gin.Context) (*usecase.DeviceProfileOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateDeviceProfileInput](ctx))
}

func (c *DeviceProfileController) GetByID(ctx *gin.Context) (*usecase.DeviceProfileOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), usecase.DeviceProfileIDInput{ID: ctx.Param("id")})
}

func (c *DeviceProfileController) List(ctx *gin.Context) ([]usecase.DeviceProfileOutput, error) {
	return c.uc.List(ctx.Request.Context(), usecase.ListDeviceProfileInput{Pagination: middleware.GetPagination(ctx)})
}

func (c *DeviceProfileController) Update(ctx *gin.Context) (*usecase.DeviceProfileOutput, error) {
	return c.uc.Update(ctx.Request.Context(), middleware.GetBodyWithID[usecase.UpdateDeviceProfileInput](ctx, "id"))
}

func (c *DeviceProfileController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), usecase.DeviceProfileIDInput{ID: ctx.Param("id")})
}
