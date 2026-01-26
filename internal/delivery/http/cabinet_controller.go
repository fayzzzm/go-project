package http

import (
	"context"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CabinetUseCase is defined by the consumer (Controller).
type CabinetUseCase interface {
	Create(ctx context.Context, input usecase.CreateCabinetInput) (*usecase.CabinetOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.CabinetOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.CabinetOutput, error)
	Update(ctx context.Context, id string, input usecase.UpdateCabinetInput) (*usecase.CabinetOutput, error)
	Delete(ctx context.Context, id string) error
}

type CabinetController struct {
	uc CabinetUseCase
}

func NewCabinetController(r gin.IRouter, uc CabinetUseCase) {
	c := &CabinetController{uc: uc}

	cabinets := r.Group("/cabinets")
	{
		cabinets.POST("", middleware.BindJSON[usecase.CreateCabinetInput](), c.Create)
		cabinets.GET("", c.List)
		cabinets.GET("/:id", c.GetByID)
		cabinets.PUT("/:id", middleware.BindJSON[usecase.UpdateCabinetInput](), c.Update)
		cabinets.DELETE("/:id", c.Delete)
	}
}

func (c *CabinetController) Create(ctx *gin.Context) {
	input := middleware.GetBody[usecase.CreateCabinetInput](ctx)

	output, err := c.uc.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, output)
}

func (c *CabinetController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	output, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, output)
}

func (c *CabinetController) List(ctx *gin.Context) {
	limit := 100
	offset := 0

	output, err := c.uc.List(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, output)
}

func (c *CabinetController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateCabinetInput](ctx)

	output, err := c.uc.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, output)
}

func (c *CabinetController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(204)
}
