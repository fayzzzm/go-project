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
	Create(ctx context.Context, input usecase.CreateCabinetInput) (*usecase.CabinetOutput, error)
	GetByID(ctx context.Context, input usecase.CabinetIDInput) (*usecase.CabinetOutput, error)
	List(ctx context.Context, input usecase.ListCabinetInput) ([]usecase.CabinetOutput, error)
	Update(ctx context.Context, input usecase.UpdateCabinetInput) (*usecase.CabinetOutput, error)
	Delete(ctx context.Context, input usecase.CabinetIDInput) error
}

type CabinetController struct {
	uc CabinetUseCase
}

func NewCabinetController(r gin.IRouter, uc CabinetUseCase) {
	c := &CabinetController{uc: uc}

	cabinets := r.Group("/cabinets")
	cabinets.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
	{
		cabinets.POST("", middleware.BindJSON[usecase.CreateCabinetInput](), utils.Handle(c.Create, http.StatusCreated))
		cabinets.GET("", utils.Handle(c.List, http.StatusOK))
		cabinets.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		cabinets.PUT("/:id", middleware.BindJSON[usecase.UpdateCabinetInput](), utils.Handle(c.Update, http.StatusOK))
		cabinets.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *CabinetController) Create(ctx *gin.Context) (*usecase.CabinetOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateCabinetInput](ctx))
}

func (c *CabinetController) GetByID(ctx *gin.Context) (*usecase.CabinetOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), usecase.CabinetIDInput{ID: ctx.Param("id")})
}

func (c *CabinetController) List(ctx *gin.Context) ([]usecase.CabinetOutput, error) {
	return c.uc.List(ctx.Request.Context(), usecase.ListCabinetInput{Pagination: middleware.GetPagination(ctx)})
}

func (c *CabinetController) Update(ctx *gin.Context) (*usecase.CabinetOutput, error) {
	return c.uc.Update(ctx.Request.Context(), middleware.GetBodyWithID[usecase.UpdateCabinetInput](ctx, "id"))
}

func (c *CabinetController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), usecase.CabinetIDInput{ID: ctx.Param("id")})
}
