package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CabinetUseCase interface {
	Create(ctx context.Context, input usecase.CreateCabinetInput, userID string) (*usecase.CabinetOutput, error)
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
	cabinets.Use(middleware.AuthMiddleware())
	{
		cabinets.POST("", middleware.BindJSON[usecase.CreateCabinetInput](), utils.Handle(c.Create, http.StatusCreated))
		cabinets.GET("", utils.Handle(c.List, http.StatusOK))
		cabinets.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		cabinets.PUT("/:id", middleware.BindJSON[usecase.UpdateCabinetInput](), utils.Handle(c.Update, http.StatusOK))
		cabinets.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *CabinetController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateCabinetInput](ctx)
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return c.uc.Create(ctx.Request.Context(), input, userID)
}

func (c *CabinetController) GetByID(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return c.uc.GetByID(ctx.Request.Context(), id)
}

func (c *CabinetController) List(ctx *gin.Context) (any, error) {
	limit, offset := middleware.GetPagination(ctx)
	return c.uc.List(ctx.Request.Context(), limit, offset)
}

func (c *CabinetController) Update(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateCabinetInput](ctx)
	return c.uc.Update(ctx.Request.Context(), id, input)
}

func (c *CabinetController) Delete(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	return nil, err
}
