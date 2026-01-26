package http

import (
	"context"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/reply"
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
		teams.POST("", middleware.BindJSON[usecase.CreateTeamInput](), c.Create)
		teams.GET("", c.List)
		teams.GET("/:id", c.GetByID)
		teams.PUT("/:id", middleware.BindJSON[usecase.UpdateTeamInput](), c.Update)
		teams.DELETE("/:id", c.Delete)
	}
}

func (c *TeamController) Create(ctx *gin.Context) {
	input := middleware.GetBody[usecase.CreateTeamInput](ctx)

	// If tenant ID is missing in body, try to get it from header/context
	if input.TenantID == "" {
		input.TenantID = ctx.GetHeader("X-Tenant-ID")
	}

	output, err := c.uc.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.Created(ctx, output)
}

func (c *TeamController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	output, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *TeamController) List(ctx *gin.Context) {
	limit := 100
	offset := 0
	tenantID := ctx.GetHeader("X-Tenant-ID")

	output, err := c.uc.List(ctx.Request.Context(), limit, offset, tenantID)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *TeamController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateTeamInput](ctx)

	output, err := c.uc.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.OK(ctx, output)
}

func (c *TeamController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.Deleted(ctx, gin.H{"id": id, "deleted": true})
}
