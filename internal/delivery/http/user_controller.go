package http

import (
	"context"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/reply"
	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	Create(ctx context.Context, input usecase.CreateUserInput) (*usecase.UserOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.UserOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.UserOutput, error)
	Update(ctx context.Context, id string, input usecase.UpdateUserInput) (*usecase.UserOutput, error)
	Delete(ctx context.Context, id string) error
}

type UserController struct {
	uc UserUseCase
}

func NewUserController(r gin.IRouter, uc UserUseCase) {
	c := &UserController{uc: uc}

	users := r.Group("/users")
	{
		users.POST("", middleware.BindJSON[usecase.CreateUserInput](), c.Create)
		users.GET("", c.List)
		users.GET("/:id", c.GetByID)
		users.PUT("/:id", middleware.BindJSON[usecase.UpdateUserInput](), c.Update)
		users.DELETE("/:id", c.Delete)
	}
}

func (c *UserController) Create(ctx *gin.Context) {
	input := middleware.GetBody[usecase.CreateUserInput](ctx)

	output, err := c.uc.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.Created(ctx, output)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	output, err := c.uc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *UserController) List(ctx *gin.Context) {
	// TODO: Parse limit/offset from query params
	limit := 100
	offset := 0

	output, err := c.uc.List(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.OK(ctx, output)
}

func (c *UserController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateUserInput](ctx)

	output, err := c.uc.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.Error(err)
		return
	}

	reply.OK(ctx, output)
}

func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	if err != nil {
		ctx.Error(err)
		return
	}
	reply.Deleted(ctx, gin.H{"id": id, "deleted": true})
}
