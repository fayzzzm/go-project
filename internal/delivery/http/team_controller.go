package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/gin-gonic/gin"
)

type TeamUseCase interface {
	Create(ctx context.Context, input usecase.CreateTeamInput) (*usecase.TeamOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.TeamOutput, error)
	List(ctx context.Context, limit, offset int, tenantID string) ([]usecase.TeamOutput, error)
	Update(ctx context.Context, id string, input usecase.UpdateTeamInput) (*usecase.TeamOutput, error)
	Delete(ctx context.Context, id string) error
}

type TeamController struct {
	uc TeamUseCase
}

func NewTeamController(r gin.IRouter, uc TeamUseCase) {
	c := &TeamController{uc: uc}

	teams := r.Group("/teams")
	{
		teams.POST("", middleware.BindJSON[usecase.CreateTeamInput](), Handle(c.Create, http.StatusCreated))
		teams.GET("", Handle(c.List, http.StatusOK))
		teams.GET("/:id", Handle(c.GetByID, http.StatusOK))
		teams.PUT("/:id", middleware.BindJSON[usecase.UpdateTeamInput](), Handle(c.Update, http.StatusOK))
		teams.DELETE("/:id", Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *TeamController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateTeamInput](ctx)
	return c.uc.Create(ctx.Request.Context(), input)
}

func (c *TeamController) GetByID(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return c.uc.GetByID(ctx.Request.Context(), id)
}

func (c *TeamController) List(ctx *gin.Context) (any, error) {
	limit, offset := middleware.GetPagination(ctx)
	tenantID := ctx.GetHeader("X-Tenant-ID")

	return c.uc.List(ctx.Request.Context(), limit, offset, tenantID)
}

func (c *TeamController) Update(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateTeamInput](ctx)
	return c.uc.Update(ctx.Request.Context(), id, input)
}

func (c *TeamController) Delete(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	return nil, err
}
